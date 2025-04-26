package services

import (
	"fluxton/models"
	"fluxton/repositories"
	"fluxton/utils"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"os"
)

const (
	ImageName = "postgrest/postgrest"
)

type PostgrestService interface {
	StartContainer(dbName string, dbPort int)
	RemoveContainer(dbName string)
}

type PostgrestServiceImpl struct {
	projectRepo *repositories.ProjectRepository
}

func NewPostgrestService(injector *do.Injector) (PostgrestService, error) {
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &PostgrestServiceImpl{
		projectRepo: projectRepo,
	}, nil
}

func (s *PostgrestServiceImpl) StartContainer(dbName string, dbPort int) {
	containerName := fmt.Sprintf("postgrest_%s", dbName)

	// Traefik label for routing based on subdomain
	traefikRule := fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s.%s`)", dbName, dbName, os.Getenv("APP_URL"))

	command := []string{
		"docker", "run", "-d", "--name", containerName,
		"--network", "fluxton_network",
		"-e", fmt.Sprintf(
			"PGRST_DB_URI=postgres://%s:%s@%s/%s",
			os.Getenv("POSTGREST_DB_USER"),
			os.Getenv("POSTGREST_DB_PASSWORD"),
			os.Getenv("POSTGREST_DB_HOST"),
			dbName,
		),
		"-e", "PGRST_DB_ANON_ROLE=" + os.Getenv("POSTGREST_DEFAULT_ROLE"),
		"-e", "PGRST_DB_SCHEMA=" + os.Getenv("POSTGREST_DEFAULT_SCHEMA"),
		"-e", "PGRST_JWT_SECRET=" + os.Getenv("JWT_SECRET"),
		"--label", "traefik.enable=true",
		"--label", traefikRule,
		"--label", fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=3000", dbName), // PostgREST default port
		ImageName,
	}

	// Execute command to start container
	if err := utils.ExecuteCommand(command); err != nil {
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
	containerName := fmt.Sprintf("postgrest_%s", dbName)

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
