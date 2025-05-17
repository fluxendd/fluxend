package services

import (
	"fluxton/constants"
	"fluxton/models"
	"fluxton/pkg"
	"fluxton/repositories"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
	"strings"
)

const (
	ImageName = "postgrest/postgrest"
)

type PostgrestConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBSchema   string
	DBRole     string
	JWTSecret  string
	AppURL     string
}

type PostgrestService interface {
	StartContainer(dbName string)
	RemoveContainer(dbName string)
	HasContainer(dbName string) bool
}

type PostgrestServiceImpl struct {
	projectRepo *repositories.ProjectRepository
	config      *PostgrestConfig
}

func NewPostgrestService(injector *do.Injector) (PostgrestService, error) {
	config := &PostgrestConfig{
		DBUser:     os.Getenv("POSTGREST_DB_USER"),
		DBPassword: os.Getenv("POSTGREST_DB_PASSWORD"),
		DBHost:     os.Getenv("POSTGREST_DB_HOST"),
		DBSchema:   os.Getenv("POSTGREST_DEFAULT_SCHEMA"),
		DBRole:     os.Getenv("POSTGREST_DEFAULT_ROLE"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AppURL:     os.Getenv("APP_URL"),
	}

	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &PostgrestServiceImpl{
		projectRepo: projectRepo,
		config:      config,
	}, nil
}

func (s *PostgrestServiceImpl) StartContainer(dbName string) {
	if err := pkg.ExecuteCommand(s.buildStartCommand(dbName)); err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to start container")

		_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, models.ProjectStatusError)
		if err != nil {
			log.Error().
				Str("action", constants.ActionPostgrest).
				Str("dbName", dbName).
				Str("error", err.Error()).
				Msg("failed to update project status to error")

			return
		}

	}

	// Update project status to active
	_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, models.ProjectStatusActive)
	if err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to update project status to active")

		return
	}

	log.Info().
		Str("action", constants.ActionPostgrest).
		Str("db", dbName).
		Msg("container started successfully")
}

func (s *PostgrestServiceImpl) RemoveContainer(dbName string) {
	containerName := s.getContainerName(dbName)

	if err := pkg.ExecuteCommand([]string{"docker", "stop", containerName}); err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to stop container")
	}

	if err := pkg.ExecuteCommand([]string{"docker", "rm", containerName}); err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to remove container")
	}

	_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, models.ProjectStatusInactive)
	if err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to update project status to inactive")

		return
	}
}

func (s *PostgrestServiceImpl) HasContainer(dbName string) bool {
	cmd := []string{"docker", "inspect", "--format='{{.State.Running}}'", s.getContainerName(dbName)}
	output, err := pkg.ExecuteCommandWithOutput(cmd)
	if err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to check if container exists")

		return false
	}

	return strings.Contains(output, "true")
}

func (s *PostgrestServiceImpl) buildStartCommand(dbName string) []string {
	return []string{
		"docker", "run", "-d", "--name", s.getContainerName(dbName),
		"--network", "fluxton_network",
		"-e", fmt.Sprintf("PGRST_DB_URI=postgres://%s:%s@%s/%s", s.config.DBUser, s.config.DBPassword, s.config.DBHost, dbName),
		"-e", "PGRST_DB_ANON_ROLE=" + s.config.DBRole,
		"-e", "PGRST_DB_SCHEMA=" + s.config.DBSchema,
		"-e", "PGRST_JWT_SECRET=" + s.config.JWTSecret,
		"--label", "traefik.enable=true",
		"--label", fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s.%s`)", dbName, dbName, s.config.AppURL),
		"--label", fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=3000", dbName),
		ImageName,
	}
}

func (s *PostgrestServiceImpl) getContainerName(dbName string) string {
	return fmt.Sprintf("postgrest_%s", dbName)
}
