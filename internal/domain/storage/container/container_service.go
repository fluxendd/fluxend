package container

import (
	"fluxton/internal/adapters/storage"
	"fluxton/internal/api/dto"
	container2 "fluxton/internal/api/dto/storage/container"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/project"
	"fluxton/internal/domain/setting"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"strings"
	"time"
)

type ContainerService interface {
	List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]Container, error)
	GetByUUID(containerUUID uuid.UUID, authUser models.AuthUser) (Container, error)
	Create(request *container2.CreateRequest, authUser models.AuthUser) (Container, error)
	Update(containerUUID uuid.UUID, authUser models.AuthUser, request *container2.CreateRequest) (*Container, error)
	Delete(request dto.DefaultRequestWithProjectHeader, containerUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type ContainerServiceImpl struct {
	settingService setting.SettingService
	projectPolicy  *project.ProjectPolicy
	containerRepo  *repositories2.ContainerRepository
	projectRepo    *repositories2.ProjectRepository
}

func NewContainerService(injector *do.Injector) (ContainerService, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*project.ProjectPolicy](injector)
	containerRepo := do.MustInvoke[*repositories2.ContainerRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)

	return &ContainerServiceImpl{
		settingService: settingService,
		projectPolicy:  policy,
		containerRepo:  containerRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *ContainerServiceImpl) List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]Container, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Container{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Container{}, errors.NewForbiddenError("container.error.listForbidden")
	}

	return s.containerRepo.ListForProject(paginationParams, projectUUID)
}

func (s *ContainerServiceImpl) GetByUUID(containerUUID uuid.UUID, authUser models.AuthUser) (Container, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return container.Container{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return container.Container{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return container.Container{}, errors.NewForbiddenError("container.error.viewForbidden")
	}

	return container, nil
}

func (s *ContainerServiceImpl) Create(request *container2.CreateRequest, authUser models.AuthUser) (Container, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return Container{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return Container{}, errors.NewForbiddenError("container.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, request.ProjectUUID)
	if err != nil {
		return Container{}, err
	}

	storageDriver := s.settingService.GetStorageDriver(request.Context)

	container := Container{
		ProjectUuid: request.ProjectUUID,
		Name:        request.Name,
		NameKey:     s.generateContainerName(),
		Provider:    storageDriver,
		IsPublic:    request.IsPublic,
		Description: request.Description,
		MaxFileSize: request.MaxFileSize,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
	}

	storageService, err := storage.GetProvider(storageDriver)
	if err != nil {
		return container.Container{}, err
	}

	createdContainer, err := storageService.CreateContainer(container.NameKey)
	if err != nil {
		return container.Container{}, err
	}

	container.Url = createdContainer

	_, err = s.containerRepo.Create(&container)
	if err != nil {
		return container.Container{}, err
	}

	return container, nil
}

func (s *ContainerServiceImpl) Update(containerUUID uuid.UUID, authUser models.AuthUser, request *container2.CreateRequest) (*Container, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return &container.Container{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &container.Container{}, errors.NewForbiddenError("container.error.updateForbidden")
	}

	err = container.PopulateModel(&container, request)
	if err != nil {
		return nil, err
	}

	container.UpdatedAt = time.Now()
	container.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, container.ProjectUuid)
	if err != nil {
		return &container.Container{}, err
	}

	return s.containerRepo.Update(&container)
}

func (s *ContainerServiceImpl) Delete(request dto.DefaultRequestWithProjectHeader, containerUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return false, err
	}

	if container.TotalFiles > 0 {
		return false, errors.NewUnprocessableError("container.error.deleteWithFiles")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("container.error.deleteForbidden")
	}

	storageService, err := storage.GetProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return false, err
	}
	err = storageService.DeleteContainer(container.NameKey)
	if err != nil {
		return false, err
	}

	return s.containerRepo.Delete(containerUUID)
}

func (s *ContainerServiceImpl) generateContainerName() string {
	containerUUID := uuid.New()

	return "container-" + strings.Replace(containerUUID.String(), "-", "", -1)
}

func (s *ContainerServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.containerRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("container.error.duplicateName")
	}

	return nil
}
