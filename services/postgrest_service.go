package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	ImageName = "postgrest/postgrest"
)

type PostgrestService interface {
	StartContainer(dbName string, dbPort int) error
	RemoveContainer(dbName string) error
}

type PostgrestServiceImpl struct{}

func NewPostgrestService() (PostgrestService, error) {
	return &PostgrestServiceImpl{}, nil
}

func (s *PostgrestServiceImpl) StartContainer(dbName string, dbPort int) error {
	containerName := fmt.Sprintf("postgrest_%s", dbName)
	command := []string{
		"docker", "run", "-d", "--name", containerName,
		"-e", fmt.Sprintf(
			"PGRST_DB_URI=postgres://%s:%s@%s/%s",
			os.Getenv("POSTGREST_DB_USER"),
			os.Getenv("POSTGREST_DB_PASSWORD"),
			os.Getenv("POSTGREST_DB_HOST"),
			dbName,
		),
		"-e", "PGRST_DB_ANON_ROLE=" + os.Getenv("POSTGREST_DEFAULT_ROLE"),
		"-e", "PGRST_DB_SCHEMA=" + os.Getenv("POSTGREST_DEFAULT_SCHEMA"),
		"-p", fmt.Sprintf("%d:3000", dbPort),
		ImageName,
	}

	if err := executeCommand(command); err != nil {
		return fmt.Errorf("failed to start container: %s", err)
	}

	return nil
}

func (s *PostgrestServiceImpl) RemoveContainer(dbName string) error {
	containerName := fmt.Sprintf("postgrest_%s", dbName)

	if err := executeCommand([]string{"docker", "stop", containerName}); err != nil {
		return fmt.Errorf("failed to stop container: %s", err)
	}

	if err := executeCommand([]string{"docker", "rm", containerName}); err != nil {
		return fmt.Errorf("failed to remove container: %s", err)
	}

	return nil
}

func executeCommand(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command failed: %s\nOutput: %s", err, string(output))

		return err
	}

	log.Printf("Command succeeded: %s", string(output))
	return nil
}
