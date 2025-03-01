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
	"strings"
	"time"
)

type BucketService interface {
	List(paginationParams utils.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Bucket, error)
	GetByUUID(bucketUUID, projectUUID uuid.UUID, authUser models.AuthUser) (models.Bucket, error)
	Create(projectUUID uuid.UUID, request *bucket_requests.CreateRequest, authUser models.AuthUser) (models.Bucket, error)
	Update(bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.CreateRequest) (*models.Bucket, error)
	Delete(bucketUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type BucketServiceImpl struct {
	s3Service     S3Service
	projectPolicy *policies.ProjectPolicy
	bucketRepo    *repositories.BucketRepository
	projectRepo   *repositories.ProjectRepository
}

func NewBucketService(injector *do.Injector) (BucketService, error) {
	s3Service, err := NewS3Service()
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	bucketRepo := do.MustInvoke[*repositories.BucketRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &BucketServiceImpl{
		s3Service:     s3Service,
		projectPolicy: policy,
		bucketRepo:    bucketRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *BucketServiceImpl) List(paginationParams utils.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Bucket, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.Bucket{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Bucket{}, errs.NewForbiddenError("bucket.error.listForbidden")
	}

	return s.bucketRepo.ListForProject(paginationParams, projectUUID)
}

func (s *BucketServiceImpl) GetByUUID(bucketUUID, projectUUID uuid.UUID, authUser models.AuthUser) (models.Bucket, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.Bucket{}, errs.NewForbiddenError("bucket.error.viewForbidden")
	}

	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	return bucket, nil
}

func (s *BucketServiceImpl) Create(projectUUID uuid.UUID, request *bucket_requests.CreateRequest, authUser models.AuthUser) (models.Bucket, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.Bucket{}, errs.NewForbiddenError("bucket.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, projectUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	bucket := models.Bucket{
		ProjectUuid: projectUUID,
		Name:        request.Name,
		AwsName:     s.generateBucketName(),
		IsPublic:    request.IsPublic,
		Description: request.Description,
		MaxFileSize: request.MaxFileSize,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
	}

	createdBucket, err := s.s3Service.CreateBucket(bucket.AwsName)
	if err != nil {
		return models.Bucket{}, err
	}

	bucket.Url = utils.PointerToString(createdBucket.Location)
	_, err = s.bucketRepo.Create(&bucket)
	if err != nil {
		return models.Bucket{}, err
	}

	return bucket, nil
}

func (s *BucketServiceImpl) Update(bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.CreateRequest) (*models.Bucket, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return &models.Bucket{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.Bucket{}, errs.NewForbiddenError("bucket.error.updateForbidden")
	}

	err = utils.PopulateModel(&bucket, request)
	if err != nil {
		return nil, err
	}

	bucket.UpdatedAt = time.Now()
	bucket.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, bucket.ProjectUuid)
	if err != nil {
		return &models.Bucket{}, err
	}

	return s.bucketRepo.Update(&bucket)
}

func (s *BucketServiceImpl) Delete(bucketUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return false, err
	}

	if bucket.TotalFiles > 0 {
		return false, errs.NewUnprocessableError("bucket.error.deleteWithFiles")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("bucket.error.deleteForbidden")
	}

	err = s.s3Service.DeleteBucket(bucket.AwsName)
	if err != nil {
		return false, err
	}

	return s.bucketRepo.Delete(bucketUUID)
}

func (s *BucketServiceImpl) generateBucketName() string {
	bucketUUID := uuid.New()

	return "bucket-" + strings.Replace(bucketUUID.String(), "-", "", -1)
}

func (s *BucketServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.bucketRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("bucket.error.duplicateName")
	}

	return nil
}
