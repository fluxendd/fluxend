package storage

import (
	"fluxton/internal/config/constants"
	"fmt"
	"github.com/labstack/echo/v4"
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
func (f *Factory) CreateProvider(ctx echo.Context, providerType string) (Provider, error) {
	switch providerType {
	case constants.StorageDriverDropbox:
		return NewDropboxService(ctx, f.injector)
	case constants.StorageDriverS3:
		return NewS3Service(ctx, f.injector)
	case constants.StorageDriverBackBlaze:
		return NewBackblazeService(ctx, f.injector)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", providerType)
	}
}
