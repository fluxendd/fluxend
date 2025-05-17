package file

import (
	storage2 "fluxton/internal/adapters/storage"
	"fluxton/internal/api/dto"
	file2 "fluxton/internal/api/dto/storage/file"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/project"
	"fluxton/internal/domain/setting"
	"fluxton/internal/domain/storage/container"
	"fluxton/models"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
	"io"
	"time"
)

type FileService interface {
	List(paginationParams dto.PaginationParams, containerUUID uuid.UUID, authUser models.AuthUser) ([]File, error)
	GetByUUID(fileUUID, containerUUID uuid.UUID, authUser models.AuthUser) (File, error)
	Create(containerUUID uuid.UUID, request *file2.CreateFileRequest, authUser models.AuthUser) (File, error)
	Rename(fileUUID, containerUUID uuid.UUID, authUser models.AuthUser, request *file2.RenameFileRequest) (*File, error)
	Delete(fileUUID, containerUUID uuid.UUID, authUser models.AuthUser, request dto.DefaultRequestWithProjectHeader) (bool, error)
}

type FileServiceImpl struct {
	settingService setting.SettingService
	projectPolicy  *project.ProjectPolicy
	containerRepo  *repositories2.ContainerRepository
	fileRepo       *repositories2.FileRepository
	projectRepo    *repositories2.ProjectRepository
}

func NewFileService(injector *do.Injector) (FileService, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*project.ProjectPolicy](injector)
	containerRepo := do.MustInvoke[*repositories2.ContainerRepository](injector)
	fileRepo := do.MustInvoke[*repositories2.FileRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)

	return &FileServiceImpl{
		settingService: settingService,
		projectPolicy:  policy,
		containerRepo:  containerRepo,
		fileRepo:       fileRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *FileServiceImpl) List(paginationParams dto.PaginationParams, containerUUID uuid.UUID, authUser models.AuthUser) ([]File, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return []File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return []File{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []File{}, errors.NewForbiddenError("file.error.listForbidden")
	}

	return s.fileRepo.ListForContainer(paginationParams, containerUUID)
}

func (s *FileServiceImpl) GetByUUID(fileUUID, containerUUID uuid.UUID, authUser models.AuthUser) (File, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return File{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return File{}, errors.NewForbiddenError("file.error.viewForbidden")
	}

	file, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return file.File{}, err
	}

	return file, nil
}

func (s *FileServiceImpl) Create(containerUUID uuid.UUID, request *file2.CreateFileRequest, authUser models.AuthUser) (File, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return File{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return File{}, errors.NewForbiddenError("file.error.createForbidden")
	}

	err = s.validate(request, container)
	if err != nil {
		return File{}, err
	}

	file := File{
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
		return file.File{}, err
	}

	storageService, err := storage2.GetProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return file.File{}, err
	}

	err = storageService.UploadFile(storage2.UploadFileInput{
		ContainerName: container.NameKey,
		FileName:      request.FullFileName,
		FileBytes:     fileBytes,
	})
	if err != nil {
		return file.File{}, err
	}

	_, err = s.fileRepo.Create(&file)
	if err != nil {
		return file.File{}, err
	}

	err = s.containerRepo.IncrementTotalFiles(containerUUID)
	if err != nil {
		return file.File{}, err
	}

	return file, nil
}

func (s *FileServiceImpl) Rename(fileUUID, containerUUID uuid.UUID, authUser models.AuthUser, request *file2.RenameFileRequest) (*File, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return nil, err
	}

	file, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return &file.File{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &file.File{}, errors.NewForbiddenError("file.error.updateForbidden")
	}

	err = s.validateNameForDuplication(request.FullFileName, container.Uuid)
	if err != nil {
		return &file.File{}, err
	}

	storageService, err := storage2.GetProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return &file.File{}, err
	}

	err = storageService.RenameFile(storage2.RenameFileInput{
		ContainerName: container.NameKey,
		FileName:      file.FullFileName,
		NewFileName:   request.FullFileName,
	})
	if err != nil {
		return nil, err
	}

	file.FullFileName = request.FullFileName
	file.UpdatedAt = time.Now()
	file.UpdatedBy = authUser.Uuid

	return s.fileRepo.Rename(&file)
}

func (s *FileServiceImpl) Delete(fileUUID, containerUUID uuid.UUID, authUser models.AuthUser, request dto.DefaultRequestWithProjectHeader) (bool, error) {
	container, err := s.containerRepo.GetByUUID(containerUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(container.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("file.error.deleteForbidden")
	}

	file, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return false, err
	}

	storageService, err := storage2.GetProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return false, err
	}

	err = storageService.DeleteFile(storage2.FileInput{
		ContainerName: container.NameKey,
		FileName:      file.FullFileName,
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

func (s *FileServiceImpl) getFileContents(request file2.CreateFileRequest) ([]byte, error) {
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

func (s *FileServiceImpl) validate(request *file2.CreateFileRequest, container container.Container) error {
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

func (s *FileServiceImpl) validateMimeType(mimeType string, container container.Container) error {
	// TODO: implement mime type validation
	// if !container.AllowedMimeTypes[mimeType] {
	//	return errors.NewUnprocessableError("file.error.invalidMimeType")
	// }

	return nil
}

func (s *FileServiceImpl) validateFileSize(fileSize int, container container.Container) error {
	if fileSize > container.MaxFileSize {
		return errors.NewUnprocessableError("file.error.sizeExceeded")
	}

	return nil
}

func (s *FileServiceImpl) validateNameForDuplication(name string, containerUUID uuid.UUID) error {
	exists, err := s.fileRepo.ExistsByNameForContainer(name, containerUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("file.error.duplicateName")
	}

	return nil
}
