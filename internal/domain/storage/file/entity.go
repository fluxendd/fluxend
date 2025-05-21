package file

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type File struct {
	shared.BaseEntity
	Uuid          uuid.UUID `db:"uuid" json:"uuid"`
	ContainerUuid uuid.UUID `db:"container_uuid" json:"containerUuid"`
	FullFileName  string    `db:"full_file_name" json:"fullFileName"`
	Size          int       `db:"size" json:"size"` // in KB
	MimeType      string    `db:"mime_type" json:"mimeType"`
	CreatedBy     uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy     uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}
