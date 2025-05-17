package storage

import (
	"fluxton/internal/config/constants"
	"fmt"
)

func GetProvider(provider string) (StorageInterface, error) {
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
