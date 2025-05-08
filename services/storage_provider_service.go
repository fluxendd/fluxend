package services

import (
	"fmt"
)

func GetStorageProvider(provider string) (StorageService, error) {
	switch provider {
	case "s3":
		return NewS3Service()
	case "dropbox":
		return NewDropboxService()
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
