package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type FileResponse struct {
	Uuid         uuid.UUID `json:"uuid"`
	BucketUuid   uuid.UUID `json:"bucketUuid"`
	FullFileName string    `json:"fullFileName"`
	Size         int       `json:"size"` // in KB
	MimeType     string    `json:"mimeType"`
	CreatedBy    uuid.UUID `json:"createdBy"`
	UpdatedBy    uuid.UUID `json:"updatedBy"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

func FileResource(file *models.File) FileResponse {
	return FileResponse{
		Uuid:         file.Uuid,
		BucketUuid:   file.BucketUuid,
		FullFileName: file.FullFileName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		CreatedBy:    file.CreatedBy,
		UpdatedBy:    file.UpdatedBy,
		CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    file.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FileResourceCollection(files []models.File) []FileResponse {
	resourceBuckets := make([]FileResponse, len(files))
	for i, file := range files {
		resourceBuckets[i] = FileResource(&file)
	}

	return resourceBuckets
}
