package postgrest

import (
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/database/client"
	"fluxton/internal/domain/project"
	"fluxton/pkg"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
	"strings"
)

const (
	ImageName = "postgrest/postgrest"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBSchema   string
	DBRole     string
	JWTSecret  string
	AppURL     string
}

type ServiceImpl struct {
	projectRepo project.Repository
	config      *Config
}

func NewPostgrestService(injector *do.Injector) (client.PostgrestService, error) {
	config := &Config{
		DBUser:     os.Getenv("POSTGREST_DB_USER"),
		DBPassword: os.Getenv("POSTGREST_DB_PASSWORD"),
		DBHost:     os.Getenv("POSTGREST_DB_HOST"),
		DBSchema:   os.Getenv("POSTGREST_DEFAULT_SCHEMA"),
		DBRole:     os.Getenv("POSTGREST_DEFAULT_ROLE"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AppURL:     os.Getenv("APP_URL"),
	}

	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ServiceImpl{
		projectRepo: projectRepo,
		config:      config,
	}, nil
}

func (s *ServiceImpl) StartContainer(dbName string) {
	if err := pkg.ExecuteCommand(s.buildStartCommand(dbName)); err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to start container")

		_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, constants.ProjectStatusError)
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
	_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, constants.ProjectStatusActive)
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

func (s *ServiceImpl) RemoveContainer(dbName string) {
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

	_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, constants.ProjectStatusInactive)
	if err != nil {
		log.Error().
			Str("action", constants.ActionPostgrest).
			Str("db", dbName).
			Str("error", err.Error()).
			Msg("failed to update project status to inactive")

		return
	}
}

func (s *ServiceImpl) HasContainer(dbName string) bool {
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

func (s *ServiceImpl) buildStartCommand(dbName string) []string {
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

func (s *ServiceImpl) getContainerName(dbName string) string {
	return fmt.Sprintf("postgrest_%s", dbName)
}
