package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/bucket_requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type FileService interface {
	List(paginationParams utils.PaginationParams, bucketUUID uuid.UUID, authUser models.AuthUser) ([]models.File, error)
	GetByUUID(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser) (models.File, error)
	Create(bucketUUID uuid.UUID, request *bucket_requests.CreateFileRequest, authUser models.AuthUser) (models.File, error)
	Rename(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.CreateFileRequest) (*models.File, error)
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

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return models.File{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.File{}, errs.NewForbiddenError("file.error.createForbidden")
	}

	// TODO: duplication check

	file := models.File{
		BucketUuid: bucketUUID,
		Name:       request.Name,
		CreatedBy:  authUser.Uuid,
		UpdatedBy:  authUser.Uuid,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdFile, err := s.s3Service.CreateBucket(bucket.AwsName)
	if err != nil {
		return models.File{}, err
	}

	utils.DumpJSON(createdFile)

	/*bucket.Url = utils.PointerToString(createdBucket.Location)
	_, err = s.fileRepo.Create(&file)
	if err != nil {
		return models.File{}, err
	}*/

	return file, nil
}

func (s *FileServiceImpl) Rename(fileUUID, bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.CreateFileRequest) (*models.File, error) {
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

	err = utils.PopulateModel(&file, request)
	if err != nil {
		return nil, err
	}

	bucket.UpdatedAt = time.Now()
	bucket.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, bucket.ProjectUuid)
	if err != nil {
		return &models.File{}, err
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

	exists, err := s.fileRepo.ExistsByUUID(fileUUID)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errs.NewNotFoundError("file.error.notFound")
	}

	err = s.s3Service.DeleteBucket(bucket.AwsName)
	if err != nil {
		return false, err
	}

	return s.bucketRepo.Delete(bucketUUID)
}

func (s *FileServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.bucketRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("bucket.error.duplicateName")
	}

	return nil
}
