package shared

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DatabaseService interface {
	Create(name string, userUUID uuid.NullUUID) error
	DropIfExists(name string) error
	Recreate(name string) error
	List() ([]string, error)
	Exists(name string) (bool, error)
	Connect(name string) (*sqlx.DB, error)
	importSeedFiles(databaseName string, userUUID uuid.UUID) error
}
