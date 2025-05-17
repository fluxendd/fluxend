package backup

import (
	"fluxton/internal/api/dto"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type Service interface {
	List(projectUUID uuid.UUID, authUser auth.AuthUser) ([]Backup, error)
	GetByUUID(backupUUID uuid.UUID, authUser auth.AuthUser) (Backup, error)
	Create(request *dto.DefaultRequestWithProjectHeader, authUser auth.AuthUser) (Backup, error)
	Delete(request dto.DefaultRequestWithProjectHeader, backupUUID uuid.UUID, authUser auth.AuthUser) (bool, error)
}

type BackupServiceImpl struct {
	projectPolicy         *project.Policy
	backupRepo            *Repository
	projectRepo           *repositories2.ProjectRepository
	backupWorkFlowService BackupWorkflowService
}

func NewBackupService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	backupRepo := do.MustInvoke[*repositories2.BackupRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)
	backupWorkFlowService := do.MustInvoke[BackupWorkflowService](injector)

	return &BackupServiceImpl{
		projectPolicy:         policy,
		backupRepo:            backupRepo,
		projectRepo:           projectRepo,
		backupWorkFlowService: backupWorkFlowService,
	}, nil
}

func (s *BackupServiceImpl) List(projectUUID uuid.UUID, authUser auth.AuthUser) ([]Backup, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Backup{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Backup{}, errors.NewForbiddenError("backup.error.listForbidden")
	}

	return s.backupRepo.ListForProject(projectUUID)
}

func (s *BackupServiceImpl) GetByUUID(backupUUID uuid.UUID, authUser auth.AuthUser) (Backup, error) {
	backup, err := s.backupRepo.GetByUUID(backupUUID)
	if err != nil {
		return Backup{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(backup.ProjectUuid)
	if err != nil {
		return Backup{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Backup{}, errors.NewForbiddenError("backup.error.viewForbidden")
	}

	return backup, nil
}

func (s *BackupServiceImpl) Create(request *dto.DefaultRequestWithProjectHeader, authUser auth.AuthUser) (Backup, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Backup{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return Backup{}, errors.NewForbiddenError("backup.error.createForbidden")
	}

	backup := Backup{
		ProjectUuid: request.ProjectUUID,
		Status:      BackupStatusCreating,
		Error:       "",
		StartedAt:   time.Now(),
	}

	createdBackup, err := s.backupRepo.Create(&backup)
	if err != nil {
		return Backup{}, err
	}

	go s.backupWorkFlowService.Create(request.Context, project.DBName, createdBackup.Uuid)

	return backup, nil
}

func (s *BackupServiceImpl) Delete(request dto.DefaultRequestWithProjectHeader, backupUUID uuid.UUID, authUser auth.AuthUser) (bool, error) {
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

	if backup.Status == BackupStatusDeleting {
		return false, errors.NewBadRequestError("backup.error.deleteInProgress")
	}

	err = s.backupRepo.UpdateStatus(backupUUID, BackupStatusDeleting, "", time.Now())
	if err != nil {
		return false, err
	}

	go s.backupWorkFlowService.Delete(request.Context, databaseName, backupUUID)

	return true, nil
}
