package file

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"mime/multipart"
)

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
