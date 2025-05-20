package storage

import (
	"github.com/guregu/null/v6"
)

type ListContainersInput struct {
	Path  string
	Limit int
	Token string
}

type ContainerMetadata struct {
	Identifier string
	Name       string
	Path       string
	Region     null.String
}

type FileInput struct {
	ContainerName string
	FileName      string
}

type UploadFileInput struct {
	ContainerName string
	FileName      string
	FileBytes     []byte
}

type RenameFileInput struct {
	ContainerName string
	FileName      string
	NewFileName   string
}
