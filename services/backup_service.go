package services

import (
	"fluxton/models"
	"fluxton/pkg/errors"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type BackupService interface {
	List(projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Backup, error)
	GetByUUID(backupUUID uuid.UUID, authUser models.AuthUser) (models.Backup, error)
	Create(request *requests.DefaultRequestWithProjectHeader, authUser models.AuthUser) (models.Backup, error)
	Delete(request requests.DefaultRequestWithProjectHeader, backupUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type BackupServiceImpl struct {
	projectPolicy         *policies.ProjectPolicy
	backupRepo            *repositories.BackupRepository
	projectRepo           *repositories.ProjectRepository
	backupWorkFlowService BackupWorkflowService
}

func NewBackupService(injector *do.Injector) (BackupService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	backupRepo := do.MustInvoke[*repositories.BackupRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	backupWorkFlowService := do.MustInvoke[BackupWorkflowService](injector)

	return &BackupServiceImpl{
		projectPolicy:         policy,
		backupRepo:            backupRepo,
		projectRepo:           projectRepo,
		backupWorkFlowService: backupWorkFlowService,
	}, nil
}

func (s *BackupServiceImpl) List(projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Backup, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.Backup{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Backup{}, errors.NewForbiddenError("backup.error.listForbidden")
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
		return models.Backup{}, errors.NewForbiddenError("backup.error.viewForbidden")
	}

	return backup, nil
}

func (s *BackupServiceImpl) Create(request *requests.DefaultRequestWithProjectHeader, authUser models.AuthUser) (models.Backup, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return models.Backup{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return models.Backup{}, errors.NewForbiddenError("backup.error.createForbidden")
	}

	backup := models.Backup{
		ProjectUuid: request.ProjectUUID,
		Status:      models.BackupStatusCreating,
		Error:       "",
		StartedAt:   time.Now(),
	}

	createdBackup, err := s.backupRepo.Create(&backup)
	if err != nil {
		return models.Backup{}, err
	}

	go s.backupWorkFlowService.Create(request.Context, project.DBName, createdBackup.Uuid)

	return backup, nil
}

func (s *BackupServiceImpl) Delete(request requests.DefaultRequestWithProjectHeader, backupUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	backup, err := s.backupRepo.GetByUUID(backupUUID)
	if err != nil {
		return false, err
	}

	databaseName, err := s.projectRepo.GetDatabaseNameByUUID(backup.ProjectUuid)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(backup.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("backup.error.deleteForbidden")
	}

	if backup.Status == models.BackupStatusDeleting {
		return false, errors.NewBadRequestError("backup.error.deleteInProgress")
	}

	err = s.backupRepo.UpdateStatus(backupUUID, models.BackupStatusDeleting, "", time.Now())
	if err != nil {
		return false, err
	}

	go s.backupWorkFlowService.Delete(request.Context, databaseName, backupUUID)

	return true, nil
}
