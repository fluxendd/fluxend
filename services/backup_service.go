package services

import (
	"fluxton/errs"
	"fluxton/models"
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
	Delete(backupUUID uuid.UUID, authUser models.AuthUser) (bool, error)
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

	go s.backupWorkFlowService.Execute(project.DBName, createdBackup.Uuid, authUser)

	return backup, nil
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
