package client

import (
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/shared"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const seedDirectory = "internal/database/seeders/client"

type Repository struct {
	db *sqlx.DB
}

func NewDatabaseRepository(injector *do.Injector) (shared.DatabaseService, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &Repository{db: db}, nil
}

func (r *Repository) Create(name string, userUUID uuid.NullUUID) error {
	_, err := r.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, name))
	if err != nil {
		log.Error().
			Str("action", constants.ActionClientDatabaseCreate).
			Str("db", name).
			Str("error", err.Error()).
			Msg("failed to create database")

		return err
	}

	if userUUID.Valid {
		return r.importSeedFiles(name, userUUID.UUID)
	}

	return nil
}

func (r *Repository) DropIfExists(name string) error {
	_, err := r.db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, name))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Recreate(name string) error {
	err := r.DropIfExists(name)
	if err != nil {
		return err
	}

	err = r.Create(name, uuid.NullUUID{})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) List() ([]string, error) {
	var databases []string
	err := r.db.Select(&databases, "SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		return []string{}, err
	}

	return databases, nil
}

func (r *Repository) Exists(name string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM pg_database WHERE datname = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Connect TODO: create actual user for using here
func (r *Repository) Connect(name string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=%s port=5432",
		os.Getenv("DATABASE_USER"),
		name,
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_SSL_MODE"),
	)

	connection, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Error().
			Str("action", constants.ActionClientDatabaseConnect).
			Str("db", name).
			Str("error", err.Error()).
			Msg("failed to connect to database")

		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	connection.DB.SetMaxOpenConns(10)
	connection.DB.SetMaxIdleConns(5)
	connection.DB.SetConnMaxLifetime(1 * time.Minute)

	return connection, nil
}

func (r *Repository) importSeedFiles(databaseName string, userUUID uuid.UUID) error {
	connection, err := r.Connect(databaseName)
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	defer connection.Close()

	// Read all files in the directory
	files, err := os.ReadDir(seedDirectory)
	if err != nil {
		log.Error().
			Str("action", constants.ActionClientDatabaseSeed).
			Str("db", databaseName).
			Str("error", err.Error()).
			Msg("Failed to read seed directory")

		return fmt.Errorf("could not read seed directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(seedDirectory, file.Name())

		// Load the contents of the SQL file
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Warn().
				Str("action", constants.ActionClientDatabaseSeed).
				Str("db", databaseName).
				Str("file", filePath).
				Str("error", err.Error()).
				Msg("Failed to read SQL file")

			continue
		}

		// Why split: If whole file is executed at once, and there is an error in one of the queries,
		// the whole file will be skipped. This way, we can execute the queries one by one.
		sqlCommands := strings.Split(string(sqlContent), ";")
		for _, query := range sqlCommands {
			if strings.Contains(query, "{{USER_ROLE}}") {
				query = strings.ReplaceAll(query, "{{USER_ROLE}}", fmt.Sprintf(`usr_%s`, strings.ReplaceAll(userUUID.String(), "-", "_")))
			}

			if _, err := connection.Exec(query); err != nil {
				log.Warn().
					Str("action", constants.ActionClientDatabaseSeed).
					Str("db", databaseName).
					Str("file", filePath).
					Str("query", query).
					Str("error", err.Error()).
					Msg("Failed to execute SQL query")

				continue
			}
		}

		log.Info().
			Str("action", constants.ActionClientDatabaseSeed).
			Str("db", databaseName).
			Str("file", filePath).
			Msg("Successfully executed SQL file")
	}

	return nil
}
