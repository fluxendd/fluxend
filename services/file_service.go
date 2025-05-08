package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/requests/bucket_requests"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
	"io"
	"time"
)

type FileService interface {
	List(paginationParams requests.PaginationParams, bucketUUID uuid.UUID, authUser models.AuthUser) ([]models.File, error)
	GetByUUID(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser) (models.File, error)
	Create(bucketUUID uuid.UUID, request *bucket_requests.CreateFileRequest, authUser models.AuthUser) (models.File, error)
	Rename(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.RenameFileRequest) (*models.File, error)
	Delete(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request requests.DefaultRequestWithProjectHeader) (bool, error)
}

type FileServiceImpl struct {
	settingService SettingService
	projectPolicy  *policies.ProjectPolicy
	bucketRepo     *repositories.BucketRepository
	fileRepo       *repositories.FileRepository
	projectRepo    *repositories.ProjectRepository
}

func NewFileService(injector *do.Injector) (FileService, error) {
	settingService, err := NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	bucketRepo := do.MustInvoke[*repositories.BucketRepository](injector)
	fileRepo := do.MustInvoke[*repositories.FileRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &FileServiceImpl{
		settingService: settingService,
		projectPolicy:  policy,
		bucketRepo:     bucketRepo,
		fileRepo:       fileRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *FileServiceImpl) List(paginationParams requests.PaginationParams, bucketUUID uuid.UUID, authUser models.AuthUser) ([]models.File, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return []models.File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return []models.File{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.File{}, errs.NewForbiddenError("file.error.listForbidden")
	}

	return s.fileRepo.ListForBucket(paginationParams, bucketUUID)
}

func (s *FileServiceImpl) GetByUUID(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser) (models.File, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return models.File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return models.File{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.File{}, errs.NewForbiddenError("file.error.viewForbidden")
	}

	file, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

func (s *FileServiceImpl) Create(bucketUUID uuid.UUID, request *bucket_requests.CreateFileRequest, authUser models.AuthUser) (models.File, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return models.File{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return models.File{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.File{}, errs.NewForbiddenError("file.error.createForbidden")
	}

	err = s.validate(request, bucket)
	if err != nil {
		return models.File{}, err
	}

	file := models.File{
		BucketUuid:   bucketUUID,
		FullFileName: request.FullFileName,
		Size:         utils.ConvertBytesToKiloBytes(int(request.File.Size)),
		MimeType:     request.File.Header.Get("Content-Type"),
		CreatedBy:    authUser.Uuid,
		UpdatedBy:    authUser.Uuid,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	fileBytes, err := s.getFileContents(*request)
	if err != nil {
		return models.File{}, err
	}

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return models.File{}, err
	}

	err = storageService.UploadFile(UploadFileInput{
		ContainerName: bucket.NameKey,
		FileName:      request.FullFileName,
		FileBytes:     fileBytes,
	})
	if err != nil {
		return models.File{}, err
	}

	_, err = s.fileRepo.Create(&file)
	if err != nil {
		return models.File{}, err
	}

	err = s.bucketRepo.IncrementTotalFiles(bucketUUID)
	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

func (s *FileServiceImpl) Rename(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.RenameFileRequest) (*models.File, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return nil, err
	}

	file, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return &models.File{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.File{}, errs.NewForbiddenError("file.error.updateForbidden")
	}

	err = s.validateNameForDuplication(request.FullFileName, bucket.Uuid)
	if err != nil {
		return &models.File{}, err
	}

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return &models.File{}, err
	}

	err = storageService.RenameFile(RenameFileInput{
		ContainerName: bucket.NameKey,
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

func (s *FileServiceImpl) Delete(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request requests.DefaultRequestWithProjectHeader) (bool, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("file.error.deleteForbidden")
	}

	file, err := s.fileRepo.GetByUUID(fileUUID)
	if err != nil {
		return false, err
	}

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return false, err
	}

	err = storageService.DeleteFile(FileInput{
		ContainerName: bucket.NameKey,
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
		err = s.bucketRepo.DecrementTotalFiles(bucketUUID)
		if err != nil {
			return false, err
		}
	}

	return fileDeleted, nil
}

func (s *FileServiceImpl) getFileContents(request bucket_requests.CreateFileRequest) ([]byte, error) {
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

func (s *FileServiceImpl) validate(request *bucket_requests.CreateFileRequest, bucket models.Bucket) error {
	fileSize := utils.ConvertBytesToKiloBytes(int(request.File.Size))

	err := s.validateMimeType(request.File.Header.Get("Content-Type"), bucket)
	if err != nil {
		return err
	}

	err = s.validateFileSize(fileSize, bucket)
	if err != nil {
		return err
	}

	err = s.validateNameForDuplication(request.FullFileName, bucket.Uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileServiceImpl) validateMimeType(mimeType string, bucket models.Bucket) error {
	// TODO: implement mime type validation
	// if !bucket.AllowedMimeTypes[mimeType] {
	//	return errs.NewUnprocessableError("file.error.invalidMimeType")
	// }

	return nil
}

func (s *FileServiceImpl) validateFileSize(fileSize int, bucket models.Bucket) error {
	if fileSize > bucket.MaxFileSize {
		return errs.NewUnprocessableError("file.error.sizeExceeded")
	}

	return nil
}

func (s *FileServiceImpl) validateNameForDuplication(name string, bucketUUID uuid.UUID) error {
	exists, err := s.fileRepo.ExistsByNameForBucket(name, bucketUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("file.error.duplicateName")
	}

	return nil
}
