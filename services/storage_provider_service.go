package services

import (
	"fmt"
)

type StorageProviderService interface {
	GetProvider(provider string) (StorageService, error)
}

type StorageProviderServiceImpl struct{}

func NewStorageProviderService() (StorageProviderService, error) {
	return &StorageProviderServiceImpl{}, nil
}

func (s *StorageProviderServiceImpl) GetProvider(provider string) (StorageService, error) {
	switch provider {
	case "s3":
		return NewS3Service()
	case "dropbox":
		return NewDropboxService()
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
