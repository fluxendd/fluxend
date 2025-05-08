package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/requests/bucket_requests"
	"github.com/google/uuid"
	"github.com/samber/do"
	"strings"
	"time"
)

type ContainerService interface {
	List(paginationParams requests.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Bucket, error)
	GetByUUID(bucketUUID uuid.UUID, authUser models.AuthUser) (models.Bucket, error)
	Create(request *bucket_requests.CreateRequest, authUser models.AuthUser) (models.Bucket, error)
	Update(bucketUUID uuid.UUID, authUser models.AuthUser, request *bucket_requests.CreateRequest) (*models.Bucket, error)
	Delete(request requests.DefaultRequestWithProjectHeader, bucketUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type BucketServiceImpl struct {
	settingService SettingService
	projectPolicy  *policies.ProjectPolicy
	bucketRepo     *repositories.BucketRepository
	projectRepo    *repositories.ProjectRepository
}

func NewContainerService(injector *do.Injector) (ContainerService, error) {
	settingService, err := NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	bucketRepo := do.MustInvoke[*repositories.BucketRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &BucketServiceImpl{
		settingService: settingService,
		projectPolicy:  policy,
		bucketRepo:     bucketRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *BucketServiceImpl) List(paginationParams requests.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Bucket, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.Bucket{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Bucket{}, errs.NewForbiddenError("bucket.error.listForbidden")
	}

	return s.bucketRepo.ListForProject(paginationParams, projectUUID)
}

func (s *BucketServiceImpl) GetByUUID(bucketUUID uuid.UUID, authUser models.AuthUser) (models.Bucket, error) {
	bucket, err := s.bucketRepo.GetByUUID(bucketUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(bucket.ProjectUuid)
	if err != nil {
		return models.Bucket{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.Bucket{}, errs.NewForbiddenError("bucket.error.viewForbidden")
	}

	return bucket, nil
}

func (s *BucketServiceImpl) Create(request *bucket_requests.CreateRequest, authUser models.AuthUser) (models.Bucket, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.Bucket{}, errs.NewForbiddenError("bucket.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, request.ProjectUUID)
	if err != nil {
		return models.Bucket{}, err
	}

	bucket := models.Bucket{
		ProjectUuid: request.ProjectUUID,
		Name:        request.Name,
		NameKey:     s.generateBucketName(),
		IsPublic:    request.IsPublic,
		Description: request.Description,
		MaxFileSize: request.MaxFileSize,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
	}

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return models.Bucket{}, err
	}

	createdContainer, err := storageService.CreateContainer(bucket.NameKey)
	if err != nil {
		return models.Bucket{}, err
	}

	bucket.Url = createdContainer

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

	err = bucket.PopulateModel(&bucket, request)
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

func (s *BucketServiceImpl) Delete(request requests.DefaultRequestWithProjectHeader, bucketUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
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

	storageService, err := GetStorageProvider(s.settingService.GetStorageDriver(request.Context))
	if err != nil {
		return false, err
	}
	err = storageService.DeleteContainer(bucket.NameKey)
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
