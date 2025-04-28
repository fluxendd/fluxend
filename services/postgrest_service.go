package services

import (
	"fluxton/models"
	"fluxton/repositories"
	"fluxton/utils"
	"fmt"
	"github.com/labstack/gommon/log"
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
	if err := utils.ExecuteCommand(s.buildStartCommand(dbName)); err != nil {
		_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, models.ProjectStatusError)
		if err != nil {
			log.Errorf("failed to update project status: %s", err)
			return
		}

		log.Errorf("failed to start container: %s", err)
	}

	// Update project status to active
	_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, models.ProjectStatusActive)
	if err != nil {
		log.Errorf("failed to update project status: %s", err)
	}
}

func (s *PostgrestServiceImpl) RemoveContainer(dbName string) {
	containerName := s.getContainerName(dbName)

	if err := utils.ExecuteCommand([]string{"docker", "stop", containerName}); err != nil {
		log.Errorf("failed to stop container: %s", err)
	}

	if err := utils.ExecuteCommand([]string{"docker", "rm", containerName}); err != nil {
		log.Errorf("failed to remove container: %s", err)
	}

	_, err := s.projectRepo.UpdateStatusByDatabaseName(dbName, models.ProjectStatusInactive)
	if err != nil {
		log.Errorf("failed to update project status: %s", err)
	}
}

func (s *PostgrestServiceImpl) HasContainer(dbName string) bool {
	cmd := []string{"docker", "inspect", "--format='{{.State.Running}}'", s.getContainerName(dbName)}
	output, err := utils.ExecuteCommandWithOutput(cmd)
	if err != nil {
		log.Errorf("failed to check container: %s", err)

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
