package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type BucketResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	ProjectUuid uuid.UUID `json:"projectUuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"isPublic"`
	Url         string    `json:"url"`
	TotalFiles  int       `json:"totalFiles"`
	MaxFileSize int       `json:"maxFileSize"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

func BucketResource(bucket *models.Bucket) BucketResponse {
	return BucketResponse{
		Uuid:        bucket.Uuid,
		ProjectUuid: bucket.ProjectUuid,
		Name:        bucket.Name,
		Description: bucket.Description,
		IsPublic:    bucket.IsPublic,
		Url:         bucket.Url,
		TotalFiles:  bucket.TotalFiles,
		MaxFileSize: bucket.MaxFileSize,
		CreatedBy:   bucket.CreatedBy,
		UpdatedBy:   bucket.UpdatedBy,
		CreatedAt:   bucket.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   bucket.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func BucketResourceCollection(buckets []models.Bucket) []BucketResponse {
	resourceBuckets := make([]BucketResponse, len(buckets))
	for i, bucket := range buckets {
		resourceBuckets[i] = BucketResource(&bucket)
	}

	return resourceBuckets
}
