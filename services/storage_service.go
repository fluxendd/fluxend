package services

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
	FileInput
	FileBytes []byte
}

type RenameFileInput struct {
	FileInput
	NewFileName string
}

type StorageService interface {
	ListContainers(input ListContainersInput) ([]string, string, error)
	CreateContainer(name string) (string, error)
	ContainerExists(name string) bool
	ShowContainer(name string) (*ContainerMetadata, error)
	DeleteContainer(name string) error
	UploadFile(input UploadFileInput) error
	RenameFile(input RenameFileInput) error
	DownloadFile(input FileInput) ([]byte, error)
	DeleteFile(input FileInput) error
}
