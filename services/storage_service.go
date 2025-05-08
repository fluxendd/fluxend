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

type StorageService interface {
	ListContainers(input ListContainersInput) ([]string, string, error)
	CreateContainer(name string) (string, error)
	ContainerExists(name string) bool
	ShowContainer(name string) (*ContainerMetadata, error)
	DeleteContainer(name string) error
}
