package backup

import (
	storage2 "fluxton/internal/adapters/storage"
	constants2 "fluxton/internal/config/constants"
	"fluxton/internal/database/repositories"
	"fluxton/internal/domain/setting"
	"fluxton/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
	"time"
)

type WorkflowService interface {
	Create(ctx echo.Context, databaseName string, backupUUID uuid.UUID)
	Delete(ctx echo.Context, databaseName string, backupUUID uuid.UUID)
}

type WorkflowServiceImpl struct {
	settingService setting.Service
	backupRepo     *repositories.BackupRepository
}

func NewBackupWorkflowService(injector *do.Injector) (WorkflowService, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	backupRepo := do.MustInvoke[*repositories.BackupRepository](injector)

	return &WorkflowServiceImpl{
		settingService: settingService,
		backupRepo:     backupRepo,
	}, nil
}

// Create pg_dump, copy file, ensure container exists, and upload to S3
func (s *WorkflowServiceImpl) Create(ctx echo.Context, databaseName string, backupUUID uuid.UUID) {
	backupFilePath := fmt.Sprintf("/tmp/%s.sql", backupUUID)

	// 1. Execute pg_dump
	if err := s.executePgDump(databaseName, backupFilePath); err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusCreatingFailed, err.Error())

		return
	}

	// 2. Copy backup file to app container
	if err := s.copyBackupToAppContainer(backupFilePath, backupUUID); err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusCreatingFailed, err.Error())

		return
	}

	// 3. Ensure backup container exists
	if err := s.ensureBackupContainerExists(ctx); err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusCreatingFailed, err.Error())

		return
	}

	fileBytes, err := s.readBackupFile(backupUUID)
	if err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusCreatingFailed, err.Error())

		return
	}

	// 4. Upload backup to S3
	err = s.uploadBackup(ctx, databaseName, backupUUID, fileBytes)
	if err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusCreatingFailed, err.Error())

		return
	}

	// 5. Update backup status to completed
	err = s.backupRepo.UpdateStatus(backupUUID, BackupStatusCreated, "", time.Now())
	if err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusCreatingFailed, err.Error())
	}

	// 6. Remove backup file from app container
	err = os.Remove(backupFilePath)
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to remove backup file from fluxton_app container")
	}
}

func (s *WorkflowServiceImpl) Delete(ctx echo.Context, databaseName string, backupUUID uuid.UUID) {
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)

	storageService, err := storage2.GetProvider(s.settingService.GetStorageDriver(ctx))
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to get storage provider")
		s.handleBackupFailure(backupUUID, BackupStatusDeletingFailed, err.Error())

		return
	}

	err = storageService.DeleteFile(storage2.FileInput{
		ContainerName: constants2.BackupContainerName,
		FileName:      filePath,
	})
	if err != nil {
		s.handleBackupFailure(backupUUID, BackupStatusDeletingFailed, err.Error())
	}

	_, err = s.backupRepo.Delete(backupUUID)
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to delete backup from database")
	}
}

func (s *WorkflowServiceImpl) executePgDump(databaseName, backupFilePath string) error {
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

	err := pkg.ExecuteCommand(command)
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupFilePath).
			Msg("failed to execute pg_dump command")

		return err
	}

	return err
}

func (s *WorkflowServiceImpl) copyBackupToAppContainer(backupFilePath string, backupUUID uuid.UUID) error {
	dockerCpCommand := []string{
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", os.Getenv("DATABASE_CONTAINER_NAME"), backupFilePath), // Source
		fmt.Sprintf("/tmp/%s.sql", backupUUID),                                     // Destination inside app container
	}

	err := pkg.ExecuteCommand(dockerCpCommand)
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("backup_uuid", backupUUID.String()).
			Msg("failed to copy backup file from fluxton_db to fluxton_app container")
	}

	return err
}

func (s *WorkflowServiceImpl) ensureBackupContainerExists(ctx echo.Context) error {
	storageService, err := storage2.GetProvider(s.settingService.GetStorageDriver(ctx))
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("error", err.Error()).
			Msg("failed to get storage provider")

		return err
	}

	containerExists := storageService.ContainerExists(constants2.BackupContainerName)
	if !containerExists {
		_, err := storageService.CreateContainer(constants2.BackupContainerName)
		if err != nil {
			log.Error().
				Str("action", constants2.ActionBackup).
				Str("container_name", constants2.BackupContainerName).
				Str("error", err.Error()).
				Msg("failed to create backup container")

			return err
		}
	}

	return nil
}

func (s *WorkflowServiceImpl) readBackupFile(backupUUID uuid.UUID) ([]byte, error) {
	fileBytes, err := os.ReadFile(fmt.Sprintf("/tmp/%s.sql", backupUUID))
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to read backup file from fluxton_app container")

		return nil, err
	}

	return fileBytes, err
}

func (s *WorkflowServiceImpl) uploadBackup(ctx echo.Context, databaseName string, backupUUID uuid.UUID, fileBytes []byte) error {
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)

	storageService, err := storage2.GetProvider(s.settingService.GetStorageDriver(ctx))
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to get storage provider")

		return err
	}

	err = storageService.UploadFile(storage2.UploadFileInput{
		ContainerName: constants2.BackupContainerName,
		FileName:      filePath,
		FileBytes:     fileBytes,
	})
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("db", databaseName).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to upload backup file to S3")

		return err
	}

	return err
}

// handleBackupFailure updates the backup status to appropriate state and logs the error
func (s *WorkflowServiceImpl) handleBackupFailure(backupUUID uuid.UUID, status, errorMessage string) {
	err := s.backupRepo.UpdateStatus(backupUUID, status, errorMessage, time.Now())
	if err != nil {
		log.Error().
			Str("action", constants2.ActionBackup).
			Str("backup_uuid", backupUUID.String()).
			Str("error", err.Error()).
			Msg("failed to update backup status in database")
		return
	}

	log.Error().
		Str("action", constants2.ActionBackup).
		Str("backup_uuid", backupUUID.String()).
		Str("status", status).
		Str("error", errorMessage).
		Msg("backup workflow failed")
}
