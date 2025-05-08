package services

import (
	"fluxton/constants"
	"fluxton/models"
	"fluxton/repositories"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
	"time"
)

type BackupWorkflowService interface {
	Create(ctx echo.Context, databaseName string, backupUUID uuid.UUID)
	Delete(ctx echo.Context, databaseName string, backupUUID uuid.UUID)
}

type BackupWorkflowServiceImpl struct {
	settingService SettingService
	backupRepo     *repositories.BackupRepository
}

func NewBackupWorkflowService(injector *do.Injector) (BackupWorkflowService, error) {
	settingService, err := NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	backupRepo := do.MustInvoke[*repositories.BackupRepository](injector)

	return &BackupWorkflowServiceImpl{
		settingService: settingService,
		backupRepo:     backupRepo,
	}, nil
}

// Create pg_dump, copy file, ensure bucket exists, and upload to S3
func (s *BackupWorkflowServiceImpl) Create(ctx echo.Context, databaseName string, backupUUID uuid.UUID) {
	backupFilePath := fmt.Sprintf("/tmp/%s.sql", backupUUID)

	// 1. Execute pg_dump
	if err := s.executePgDump(databaseName, backupFilePath); err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 2. Copy backup file to app container
	if err := s.copyBackupToAppContainer(backupFilePath, backupUUID); err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 3. Ensure backup bucket exists
	if err := s.ensureBackupContainerExists(ctx); err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	fileBytes, err := s.readBackupFile(backupUUID)
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 4. Upload backup to S3
	err = s.uploadBackup(ctx, databaseName, backupUUID, fileBytes)
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 5. Update backup status to completed
	err = s.backupRepo.UpdateStatus(backupUUID, models.BackupStatusCreated, "", time.Now())
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())
	}

	// 6. Remove backup file from app container
	err = os.Remove(backupFilePath)
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to remove backup file from fluxton_app container")
	}
}

func (s *BackupWorkflowServiceImpl) Delete(ctx echo.Context, databaseName string, backupUUID uuid.UUID) {
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(ctx))
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to get storage provider")
		s.handleBackupFailure(backupUUID, models.BackupStatusDeletingFailed, err.Error())

		return
	}

	err = storageService.DeleteFile(FileInput{
		ContainerName: constants.BackupBucketName,
		FileName:      filePath,
	})
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusDeletingFailed, err.Error())
	}

	_, err = s.backupRepo.Delete(backupUUID)
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to delete backup from database")
	}
}

func (s *BackupWorkflowServiceImpl) executePgDump(databaseName, backupFilePath string) error {
	command := []string{
		"docker",
		"exec",
		os.Getenv("DATABASE_CONTAINER_NAME"),
		"pg_dump",
		"-U",
		os.Getenv("DATABASE_USER"),
		"-d",
		databaseName,
		"-f", backupFilePath,
	}

	err := utils.ExecuteCommand(command)
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupFilePath).
			Msg("failed to execute pg_dump command")

		return err
	}

	return err
}

func (s *BackupWorkflowServiceImpl) copyBackupToAppContainer(backupFilePath string, backupUUID uuid.UUID) error {
	dockerCpCommand := []string{
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", os.Getenv("DATABASE_CONTAINER_NAME"), backupFilePath), // Source
		fmt.Sprintf("/tmp/%s.sql", backupUUID),                                     // Destination inside app container
	}

	err := utils.ExecuteCommand(dockerCpCommand)
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("backup_uuid", backupUUID.String()).
			Msg("failed to copy backup file from fluxton_db to fluxton_app container")
	}

	return err
}

func (s *BackupWorkflowServiceImpl) ensureBackupContainerExists(ctx echo.Context) error {
	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(ctx))
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("error", err.Error()).
			Msg("failed to get storage provider")

		return err
	}

	bucketExists := storageService.ContainerExists(constants.BackupBucketName)
	if !bucketExists {
		_, err := storageService.CreateContainer(constants.BackupBucketName)
		if err != nil {
			log.Error().
				Str("action", constants.ActionBackup).
				Str("bucket_name", constants.BackupBucketName).
				Str("error", err.Error()).
				Msg("failed to create backup bucket")

			return err
		}
	}

	return nil
}

func (s *BackupWorkflowServiceImpl) readBackupFile(backupUUID uuid.UUID) ([]byte, error) {
	fileBytes, err := os.ReadFile(fmt.Sprintf("/tmp/%s.sql", backupUUID))
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to read backup file from fluxton_app container")

		return nil, err
	}

	return fileBytes, err
}

func (s *BackupWorkflowServiceImpl) uploadBackup(ctx echo.Context, databaseName string, backupUUID uuid.UUID, fileBytes []byte) error {
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(ctx))
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to get storage provider")

		return err
	}

	err = storageService.UploadFile(UploadFileInput{
		ContainerName: constants.BackupBucketName,
		FileName:      filePath,
		FileBytes:     fileBytes,
	})
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to upload backup file to S3")

		return err
	}

	return err
}

// handleBackupFailure updates the backup status to appropriate state and logs the error
func (s *BackupWorkflowServiceImpl) handleBackupFailure(backupUUID uuid.UUID, status, errorMessage string) {
	err := s.backupRepo.UpdateStatus(backupUUID, status, errorMessage, time.Now())
	if err != nil {
		log.Error().
			Str("action", constants.ActionBackup).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to update backup status in database")
		return
	}

	log.Error().
		Str("action", constants.ActionBackup).
		Str("backup_uuid", backupUUID.String()).
		Str("status", status).
		Str("error", errorMessage).
		Msg("backup workflow failed")
}
