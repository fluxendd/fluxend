package file

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"time"
)

type File struct {
	shared.BaseModel
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

type CreateFileInput struct {
	Context      echo.Context
	ProjectUUID  uuid.UUID             `db:"project_uuid" json:"projectUUID"`
	FullFileName string                `json:"-" form:"full_file_name"`
	File         *multipart.FileHeader `json:"-" form:"file"`
}

type RenameFileInput struct {
	Context      echo.Context
	ProjectUUID  uuid.UUID `db:"project_uuid" json:"projectUUID"`
	FullFileName string    `json:"full_file_name"`
}

func (u File) GetTableName() string {
	return "storage.files"
}
