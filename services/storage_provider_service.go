package services

import (
	"fluxton/constants"
	"fmt"
)

func GetStorageProvider(provider string) (StorageService, error) {
	switch provider {
	case constants.StorageDriverS3:
		return NewS3Service()
	case constants.StorageDriverDropbox:
		return NewDropboxService()
	case constants.StorageDriverBackBlaze:
		return NewBackblazeService()
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
