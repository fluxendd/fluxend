package storage

import (
	"fluxton/internal/config/constants"
	"fmt"
	"github.com/samber/do"
)

type Provider interface {
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

type Factory struct {
	injector *do.Injector
}

func NewFactory(injector *do.Injector) (*Factory, error) {
	return &Factory{injector: injector}, nil
}

// TODO: use injector and context
func (f *Factory) CreateProvider(providerType string) (Provider, error) {
	switch providerType {
	case constants.StorageDriverDropbox:
		return NewDropboxProvider(f.injector)
	case constants.StorageDriverS3:
		return NewS3Provider(f.injector)
	case constants.StorageDriverBackBlaze:
		return NewBackblazeProvider(f.injector)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", providerType)
	}
}
