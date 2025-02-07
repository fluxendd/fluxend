package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseRepository struct {
	db *sqlx.DB
}

func NewDatabaseRepository(injector *do.Injector) (*DatabaseRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &DatabaseRepository{db: db}, nil
}

func (r *DatabaseRepository) Create(name string) error {
	_, err := r.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, name))
	if err != nil {
		return err
	}

	return r.importSeedFiles(name)
}

func (r *DatabaseRepository) DropIfExists(name string) error {
	_, err := r.db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, name))
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseRepository) Recreate(name string) error {
	err := r.DropIfExists(name)
	if err != nil {
		return err
	}

	err = r.Create(name)
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseRepository) List() ([]string, error) {
	var databases []string
	err := r.db.Select(&databases, "SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		return []string{}, err
	}

	return databases, nil
}

func (r *DatabaseRepository) Exists(name string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM pg_database WHERE datname = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *DatabaseRepository) Connect(name string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf("dbname=%s sslmode=disable", name))
}

func (r *DatabaseRepository) importSeedFiles(databaseName string) error {
	connection, err := r.Connect(databaseName)
	defer connection.Close()

	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}

	seedDir := "seeders/client"

	// Read all files in the directory
	files, err := os.ReadDir(seedDir)
	if err != nil {
		return fmt.Errorf("could not read seed directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(seedDir, file.Name())

		// Load the contents of the SQL file
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Debug("DB: %s => Skipping file %s: could not read file: %v", databaseName, filePath, err)

			continue
		}

		// Execute the SQL statements
		if _, err := connection.Exec(string(sqlContent)); err != nil {
			log.Debug("DB: %s => Skipping file %s: could not execute SQL: %v", databaseName, filePath, err)

			continue
		}

		log.Debug("DB: %s => Successfully executed seed file %s", databaseName, filePath)
	}

	return nil
}
