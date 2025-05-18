package container

import (
	"fluxton/internal/adapters/storage"
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/storage/container"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/internal/domain/setting"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"strings"
	"time"
)

type Service interface {
	List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser auth.User) ([]Container, error)
	GetByUUID(containerUUID uuid.UUID, authUser auth.User) (Container, error)
	Create(request *container.CreateRequest, authUser auth.User) (Container, error)
	Update(containerUUID uuid.UUID, authUser auth.User, request *container.CreateRequest) (*Container, error)
	Delete(request dto.DefaultRequestWithProjectHeader, containerUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ServiceImpl struct {
	settingService setting.Service
	projectPolicy  *project.Policy
	containerRepo  Repository
	projectRepo    project.Repository
}

func NewContainerService(injector *do.Injector) (Service, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*project.Policy](injector)
	containerRepo := do.MustInvoke[Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ServiceImpl{
		settingService: settingService,
		projectPolicy:  policy,
		containerRepo:  containerRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *ServiceImpl) List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser auth.User) ([]Container, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Container{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Container{}, errors.NewForbiddenError("container.error.listForbidden")
	}

	return s.containerRepo.ListForProject(paginationParams, projectUUID)
}

func (s *ServiceImpl) GetByUUID(containerUUID uuid.UUID, authUser auth.User) (Container, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return Container{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return Container{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Container{}, errors.NewForbiddenError("container.error.viewForbidden")
	}

	return fetchedContainer, nil
}

func (s *ServiceImpl) Create(request *container.CreateRequest, authUser auth.User) (Container, error) {
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

	containerInput := Container{
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
		return Container{}, err
	}

	createdContainer, err := storageService.CreateContainer(containerInput.NameKey)
	if err != nil {
		return Container{}, err
	}

	containerInput.Url = createdContainer

	_, err = s.containerRepo.Create(&containerInput)
	if err != nil {
		return Container{}, err
	}

	return containerInput, nil
}

func (s *ServiceImpl) Update(containerUUID uuid.UUID, authUser auth.User, request *container.CreateRequest) (*Container, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return &Container{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &Container{}, errors.NewForbiddenError("container.error.updateForbidden")
	}

	/*err = container.PopulateModel(&fetchedContainer, request)
	if err != nil {
		return nil, err
	}*/

	fetchedContainer.UpdatedAt = time.Now()
	fetchedContainer.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, fetchedContainer.ProjectUuid)
	if err != nil {
		return &Container{}, err
	}

	return s.containerRepo.Update(&fetchedContainer)
}

func (s *ServiceImpl) Delete(request dto.DefaultRequestWithProjectHeader, containerUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return false, err
	}

	if fetchedContainer.TotalFiles > 0 {
		return false, errors.NewUnprocessableError("container.error.deleteWithFiles")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
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
	err = storageService.DeleteContainer(fetchedContainer.NameKey)
	if err != nil {
		return false, err
	}

	return s.containerRepo.Delete(containerUUID)
}

func (s *ServiceImpl) generateContainerName() string {
	containerUUID := uuid.New()

	return "container-" + strings.Replace(containerUUID.String(), "-", "", -1)
}

func (s *ServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.containerRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("container.error.duplicateName")
	}

	return nil
}
