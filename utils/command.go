package utils

import (
	"github.com/labstack/gommon/log"
	"os/exec"
)

func ExecuteCommand(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Command failed: %s\nOutput: %s", err, string(output))

		return err
	}

	log.Printf("Command succeeded: cmd=[%s] : output=[%s]", command, string(output))

	return nil
}

func ExecuteCommandWithOutput(command []string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
