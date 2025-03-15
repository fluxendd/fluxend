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
	backup.CompletedAt = time.Now()

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
	command := []string{
		"pg_dump",
		"-U", os.Getenv("POSTGRES_USER"),
		"-d", databaseName,
		"-f", fmt.Sprintf("/tmp/%s.sql", backupUUID),
	}

	err := utils.ExecuteCommand(command)
	if err != nil {
		_, err := s.Update(backupUUID, models.BackupStatusFailed, err.Error(), authUser)
		if err != nil {
			log.Errorf("failed to update backup status: %s", err)
		}
	}

	fileBytes, err := os.ReadFile(fmt.Sprintf("/tmp/%s.sql", backupUUID))
	if err != nil {
		log.Errorf("failed to read backup file: %s", err)
	}

	filePath := fmt.Sprintf("%s/%s/%s.sql", configs.BackupBucketName, databaseName, backupUUID)
	err = s.s3Service.UploadFile(configs.BackupBucketName, filePath, fileBytes)
	if err != nil {
		_, err := s.Update(backupUUID, models.BackupStatusFailed, err.Error(), authUser)
		if err != nil {
			log.Errorf("failed to update backup status: %s", err)
		}
	}

	_, err = s.Update(backupUUID, models.BackupStatusCompleted, "", authUser)
	if err != nil {
		log.Errorf("failed to update backup status: %s", err)
	}
}
