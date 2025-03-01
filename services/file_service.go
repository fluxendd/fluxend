package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/bucket_requests"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
	"io"
	"time"
)

type FileService interface {
	List(paginationParams utils.PaginationParams, bucketUUID uuid.UUID, authUser models.AuthUser) ([]models.File, error)
	GetByUUID(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser) (models.File, error)
	Create(bucketUUID uuid.UUID, request *bucket_requests.CreateFileRequest, authUser models.AuthUser) (models.File, error)
	Rename(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.RenameFileRequest) (*models.File, error)
	Delete(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FileServiceImpl struct {
	s3Service     S3Service
	projectPolicy *policies.ProjectPolicy
	bucketRepo    *repositories.BucketRepository
	fileRepo      *repositories.FileRepository
	projectRepo   *repositories.ProjectRepository
}

func NewFileService(injector *do.Injector) (FileService, error) {
	s3Service, err := NewS3Service()
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	bucketRepo := do.MustInvoke[*repositories.BucketRepository](injector)
	fileRepo := do.MustInvoke[*repositories.FileRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &FileServiceImpl{
		s3Service:     s3Service,
		projectPolicy: policy,
		bucketRepo:    bucketRepo,
		fileRepo:      fileRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FileServiceImpl) List(paginationParams utils.PaginationParams, bucketUUID uuid.UUID, authUser models.AuthUser) ([]models.File, error) {
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

	fileSize := utils.BytesToKiloBytes(int(request.File.Size))

	if fileSize > bucket.MaxFileSize {
		return models.File{}, errs.NewUnprocessableError("file.error.sizeExceeded")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return models.File{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.File{}, errs.NewForbiddenError("file.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, bucketUUID)
	if err != nil {
		return models.File{}, err
	}

	file := models.File{
		BucketUuid: bucketUUID,
		Name:       request.Name,
		Size:       fileSize,
		MimeType:   request.File.Header.Get("Content-Type"),
		Path:       request.Name,
		CreatedBy:  authUser.Uuid,
		UpdatedBy:  authUser.Uuid,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fileBytes, err := s.getFileContents(*request)
	if err != nil {
		return models.File{}, err
	}

	err = s.s3Service.UploadFile(bucket.AwsName, request.Name, fileBytes)
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

	file.Name = request.Name
	file.Path = request.Name
	file.UpdatedAt = time.Now()
	file.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, bucket.ProjectUuid)
	if err != nil {
		return &models.File{}, err
	}

	err = s.s3Service.RenameFile(bucket.AwsName, file.Path, request.Name)
	if err != nil {
		return nil, err
	}

	return s.fileRepo.Rename(&file)
}

func (s *FileServiceImpl) Delete(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
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

	err = s.s3Service.DeleteFile(bucket.AwsName, file.Path)
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

func (s *FileServiceImpl) validateNameForDuplication(name string, fileUUID uuid.UUID) error {
	exists, err := s.fileRepo.ExistsByNameForBucket(name, fileUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("file.error.duplicateName")
	}

	return nil
}
