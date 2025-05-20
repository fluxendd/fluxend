package file

import (
	"fluxton/internal/adapters/storage"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/internal/domain/setting"
	"fluxton/internal/domain/shared"
	"fluxton/internal/domain/storage/container"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"io"
	"time"
)

type Service interface {
	List(paginationParams shared.PaginationParams, containerUUID uuid.UUID, authUser auth.User) ([]File, error)
	GetByUUID(fileUUID, containerUUID uuid.UUID, authUser auth.User) (File, error)
	Create(containerUUID uuid.UUID, request *CreateFileInput, authUser auth.User) (File, error)
	Rename(fileUUID, containerUUID uuid.UUID, authUser auth.User, request *RenameFileInput) (*File, error)
	Delete(fileUUID, containerUUID uuid.UUID, authUser auth.User, context echo.Context) (bool, error)
}

type ServiceImpl struct {
	settingService setting.Service
	projectPolicy  *project.Policy
	containerRepo  container.Repository
	fileRepo       Repository
	projectRepo    project.Repository
	storageFactory *storage.Factory
}

func NewFileService(injector *do.Injector) (Service, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*project.Policy](injector)
	containerRepo := do.MustInvoke[container.Repository](injector)
	fileRepo := do.MustInvoke[Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)
	storageFactory := do.MustInvoke[*storage.Factory](injector)

	return &ServiceImpl{
		settingService: settingService,
		projectPolicy:  policy,
		containerRepo:  containerRepo,
		fileRepo:       fileRepo,
		projectRepo:    projectRepo,
		storageFactory: storageFactory,
	}, nil
}

func (s *ServiceImpl) List(paginationParams shared.PaginationParams, containerUUID uuid.UUID, authUser auth.User) ([]File, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return []File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return []File{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []File{}, errors.NewForbiddenError("file.error.listForbidden")
	}

	return s.fileRepo.ListForContainer(paginationParams, containerUUID)
}

func (s *ServiceImpl) GetByUUID(fileUUID, containerUUID uuid.UUID, authUser auth.User) (File, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return File{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return File{}, errors.NewForbiddenError("file.error.viewForbidden")
	}

	fetchedFile, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return File{}, err
	}

	return fetchedFile, nil
}

func (s *ServiceImpl) Create(containerUUID uuid.UUID, request *CreateFileInput, authUser auth.User) (File, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return File{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return File{}, errors.NewForbiddenError("file.error.createForbidden")
	}

	err = s.validate(request, fetchedContainer)
	if err != nil {
		return File{}, err
	}

	fileInput := File{
		ContainerUuid: containerUUID,
		FullFileName:  request.FullFileName,
		Size:          pkg.ConvertBytesToKiloBytes(int(request.File.Size)),
		MimeType:      request.File.Header.Get("Content-Type"),
		CreatedBy:     authUser.Uuid,
		UpdatedBy:     authUser.Uuid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	fileBytes, err := s.getFileContents(*request)
	if err != nil {
		return File{}, err
	}

	storageService, err := s.storageFactory.CreateProvider(request.Context, s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return File{}, err
	}

	err = storageService.UploadFile(storage.UploadFileInput{
		ContainerName: fetchedContainer.NameKey,
		FileName:      request.FullFileName,
		FileBytes:     fileBytes,
	})
	if err != nil {
		return File{}, err
	}

	_, err = s.fileRepo.Create(&fileInput)
	if err != nil {
		return File{}, err
	}

	err = s.containerRepo.IncrementTotalFiles(containerUUID)
	if err != nil {
		return File{}, err
	}

	return fileInput, nil
}

func (s *ServiceImpl) Rename(fileUUID, containerUUID uuid.UUID, authUser auth.User, request *RenameFileInput) (*File, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return nil, err
	}

	fetchedFile, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return &File{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &File{}, errors.NewForbiddenError("file.error.updateForbidden")
	}

	err = s.validateNameForDuplication(request.FullFileName, fetchedContainer.Uuid)
	if err != nil {
		return &File{}, err
	}

	storageService, err := s.storageFactory.CreateProvider(request.Context, s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return &File{}, err
	}

	err = storageService.RenameFile(storage.RenameFileInput{
		ContainerName: fetchedContainer.NameKey,
		FileName:      fetchedFile.FullFileName,
		NewFileName:   request.FullFileName,
	})
	if err != nil {
		return nil, err
	}

	fetchedFile.FullFileName = request.FullFileName
	fetchedFile.UpdatedAt = time.Now()
	fetchedFile.UpdatedBy = authUser.Uuid

	return s.fileRepo.Rename(&fetchedFile)
}

func (s *ServiceImpl) Delete(fileUUID, containerUUID uuid.UUID, authUser auth.User, context echo.Context) (bool, error) {
	fetchedContainer, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedContainer.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("file.error.deleteForbidden")
	}

	fetchedFile, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return false, err
	}

	storageService, err := s.storageFactory.CreateProvider(context, s.settingService.GetStorageDriver(context))
	if err != nil {
		return false, err
	}

	err = storageService.DeleteFile(storage.FileInput{
		ContainerName: fetchedContainer.NameKey,
		FileName:      fetchedFile.FullFileName,
	})
	if err != nil {
		return false, err
	}

	fileDeleted, err := s.fileRepo.Delete(fileUUID)
	if err != nil {
		return false, err
	}

	if fileDeleted {
		err = s.containerRepo.DecrementTotalFiles(containerUUID)
		if err != nil {
			return false, err
		}
	}

	return fileDeleted, nil
}

func (s *ServiceImpl) getFileContents(request CreateFileInput) ([]byte, error) {
	fileHandler, err := request.File.Open() // Open the file
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileHandler.Close()

	fileBytes, err := io.ReadAll(fileHandler) // Read the file content as bytes
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return fileBytes, nil
}

func (s *ServiceImpl) validate(request *CreateFileInput, container container.Container) error {
	fileSize := pkg.ConvertBytesToKiloBytes(int(request.File.Size))

	err := s.validateMimeType(request.File.Header.Get("Content-Type"), container)
	if err != nil {
		return err
	}

	err = s.validateFileSize(fileSize, container)
	if err != nil {
		return err
	}

	err = s.validateNameForDuplication(request.FullFileName, container.Uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) validateMimeType(mimeType string, container container.Container) error {
	// TODO: implement mime type validation
	// if !container.AllowedMimeTypes[mimeType] {
	//	return errors.NewUnprocessableError("file.error.invalidMimeType")
	// }

	return nil
}

func (s *ServiceImpl) validateFileSize(fileSize int, container container.Container) error {
	if fileSize > container.MaxFileSize {
		return errors.NewUnprocessableError("file.error.sizeExceeded")
	}

	return nil
}

func (s *ServiceImpl) validateNameForDuplication(name string, containerUUID uuid.UUID) error {
	exists, err := s.fileRepo.ExistsByNameForContainer(name, containerUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("file.error.duplicateName")
	}

	return nil
}
