package file

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid          uuid.UUID `json:"uuid"`
	ContainerUuid uuid.UUID `json:"containerUuid"`
	FullFileName  string    `json:"fullFileName"`
	Size          int       `json:"size"` // in KB
	MimeType      string    `json:"mimeType"`
	CreatedBy     uuid.UUID `json:"createdBy"`
	UpdatedBy     uuid.UUID `json:"updatedBy"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
}

type DownloadResponse struct {
	Url       string `json:"url"`
	ExpiresIn int64  `json:"expiresIn"` // in seconds
}
