package services

import (
	"fluxton/configs"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"os"
	"time"
)

type BackupService interface {
	List(projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Backup, error)
	GetByUUID(backupUUID uuid.UUID, authUser models.AuthUser) (models.Backup, error)
	Create(request *requests.DefaultRequestWithProjectHeader, authUser models.AuthUser) (models.Backup, error)
	Update(backupUUID uuid.UUID, status, error string, authUser models.AuthUser) (*models.Backup, error)
	Delete(backupUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type BackupServiceImpl struct {
	s3Service     S3Service
	projectPolicy *policies.ProjectPolicy
	backupRepo    *repositories.BackupRepository
	projectRepo   *repositories.ProjectRepository
}

func NewBackupService(injector *do.Injector) (BackupService, error) {
	s3Service, err := NewS3Service()
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	backupRepo := do.MustInvoke[*repositories.BackupRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &BackupServiceImpl{
		s3Service:     s3Service,
		projectPolicy: policy,
		backupRepo:    backupRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *BackupServiceImpl) List(projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Backup, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.Backup{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Backup{}, errs.NewForbiddenError("backup.error.listForbidden")
	}

	return s.backupRepo.ListForProject(projectUUID)
}

func (s *BackupServiceImpl) GetByUUID(backupUUID uuid.UUID, authUser models.AuthUser) (models.Backup, error) {
	backup, err := s.backupRepo.GetByUUID(backupUUID)
	if err != nil {
		return models.Backup{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(backup.ProjectUuid)
	if err != nil {
		return models.Backup{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.Backup{}, errs.NewForbiddenError("backup.error.viewForbidden")
	}

	return backup, nil
}

func (s *BackupServiceImpl) Create(request *requests.DefaultRequestWithProjectHeader, authUser models.AuthUser) (models.Backup, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return models.Backup{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return models.Backup{}, errs.NewForbiddenError("backup.error.createForbidden")
	}

	backup := models.Backup{
		ProjectUuid: request.ProjectUUID,
		Status:      models.BackupStatusPending,
		Error:       "",
		StartedAt:   time.Now(),
	}

	createdBackup, err := s.backupRepo.Create(&backup)
	if err != nil {
		return models.Backup{}, err
	}

	go s.SaveBackup(project.DBName, createdBackup.Uuid, authUser)

	return backup, nil
}

func (s *BackupServiceImpl) Update(backupUUID uuid.UUID, status, error string, authUser models.AuthUser) (*models.Backup, error) {
	backup, err := s.backupRepo.GetByUUID(backupUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(backup.ProjectUuid)
	if err != nil {
		return &models.Backup{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.Backup{}, errs.NewForbiddenError("backup.error.updateForbidden")
	}

	backup.Status = status
	backup.Error = error
	now := time.Now()
	backup.CompletedAt = &now

	return s.backupRepo.Update(&backup)
}

func (s *BackupServiceImpl) Delete(backupUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	backup, err := s.backupRepo.GetByUUID(backupUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(backup.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("backup.error.deleteForbidden")
	}

	return s.backupRepo.Delete(backupUUID)
}

// SaveBackup dumps the backup to the storage and updates status, runs as a goroutine
func (s *BackupServiceImpl) SaveBackup(databaseName string, backupUUID uuid.UUID, authUser models.AuthUser) {
	// 1. Execute pg_dump in PostgreSQL container to take the backup
	backupFilePath := fmt.Sprintf("/tmp/%s.sql", backupUUID)

	command := []string{
		"docker",
		"exec",
		os.Getenv("DATABASE_CONTAINER_NAME"), // e.g. "fluxton_db"
		"pg_dump",
		"-U",
		os.Getenv("DATABASE_USER"),
		"-d",
		databaseName,
		"-f", backupFilePath, // Directly specify the output file
	}

	err := utils.ExecuteCommand(command)
	if err != nil {
		// Update backup status to failed if pg_dump command fails
		_, statusUpdateErr := s.Update(backupUUID, models.BackupStatusFailed, err.Error(), authUser)
		if statusUpdateErr != nil {
			log.Errorf("failed to update backup status: %s", statusUpdateErr.Error())
		}

		return
	}

	// 2. Copy the backup file from the PostgreSQL container to the app container (fluxton_app)
	dockerCpCommand := []string{
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", os.Getenv("DATABASE_CONTAINER_NAME"), backupFilePath), // Source
		fmt.Sprintf("/tmp/%s.sql", backupUUID),                                     // Destination inside fluxton_app container
	}

	cmdError := utils.ExecuteCommand(dockerCpCommand)
	if cmdError != nil {
		// Update backup status to failed if docker cp command fails
		_, statusUpdateErr := s.Update(backupUUID, models.BackupStatusFailed, cmdError.Error(), authUser)
		if statusUpdateErr != nil {
			log.Errorf("failed to update backup status: %s", statusUpdateErr.Error())
		}

		log.Errorf("failed to copy backup file: %s", cmdError.Error())
		return
	}

	// 3. Check if backup bucket exists, if not create it
	bucketExists := s.s3Service.BucketExists(configs.BackupBucketName)
	if !bucketExists {
		_, err = s.s3Service.CreateBucket(configs.BackupBucketName)
		if err != nil {
			// Update backup status to failed if creating bucket fails
			_, statusUpdateErr := s.Update(backupUUID, models.BackupStatusFailed, err.Error(), authUser)
			if statusUpdateErr != nil {
				log.Errorf("failed to update backup status: %s", statusUpdateErr.Error())
			}

			log.Errorf("failed to create backup bucket: %s", err.Error())
			return
		}
	}

	// 4. Read the backup file from the app container
	fileBytes, fileReadErr := os.ReadFile(fmt.Sprintf("/tmp/%s.sql", backupUUID))
	if fileReadErr != nil {
		// Update backup status to failed if reading the file fails
		_, statusUpdateErr := s.Update(backupUUID, models.BackupStatusFailed, fileReadErr.Error(), authUser)
		if statusUpdateErr != nil {
			log.Errorf("failed to update backup status: %s", statusUpdateErr.Error())
		}

		log.Errorf("failed to read backup file: %s", fileReadErr.Error())
		return
	}

	// 5. Upload the backup file to S3
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)
	uploadErr := s.s3Service.UploadFile(configs.BackupBucketName, filePath, fileBytes)
	if uploadErr != nil {
		// Update backup status to failed if S3 upload fails
		_, statusUpdateErr := s.Update(backupUUID, models.BackupStatusFailed, uploadErr.Error(), authUser)
		if statusUpdateErr != nil {
			log.Errorf("failed to update backup status: %s", err)
		}

		log.Errorf("failed to upload backup to S3: %s", uploadErr.Error())

		return
	}

	// 6. Update backup status to completed
	_, err = s.Update(backupUUID, models.BackupStatusCompleted, "", authUser)
	if err != nil {
		log.Errorf("failed to update backup status: %s", err)
	}
}
