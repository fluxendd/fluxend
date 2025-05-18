package backup

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type Service interface {
	List(projectUUID uuid.UUID, authUser auth.User) ([]Backup, error)
	GetByUUID(backupUUID uuid.UUID, authUser auth.User) (Backup, error)
	Create(request *dto.DefaultRequestWithProjectHeader, authUser auth.User) (Backup, error)
	Delete(request dto.DefaultRequestWithProjectHeader, backupUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ServiceImpl struct {
	projectPolicy         *project.Policy
	backupRepo            Repository
	projectRepo           project.Repository
	backupWorkFlowService WorkflowService
}

func NewBackupService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	backupRepo := do.MustInvoke[Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)
	backupWorkFlowService := do.MustInvoke[WorkflowService](injector)

	return &ServiceImpl{
		projectPolicy:         policy,
		backupRepo:            backupRepo,
		projectRepo:           projectRepo,
		backupWorkFlowService: backupWorkFlowService,
	}, nil
}

func (s *ServiceImpl) List(projectUUID uuid.UUID, authUser auth.User) ([]Backup, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Backup{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Backup{}, errors.NewForbiddenError("backup.error.listForbidden")
	}

	return s.backupRepo.ListForProject(projectUUID)
}

func (s *ServiceImpl) GetByUUID(backupUUID uuid.UUID, authUser auth.User) (Backup, error) {
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

func (s *ServiceImpl) Create(request *dto.DefaultRequestWithProjectHeader, authUser auth.User) (Backup, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Backup{}, err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
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

	go s.backupWorkFlowService.Create(request.Context, fetchedProject.DBName, createdBackup.Uuid)

	return backup, nil
}

func (s *ServiceImpl) Delete(request dto.DefaultRequestWithProjectHeader, backupUUID uuid.UUID, authUser auth.User) (bool, error) {
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
