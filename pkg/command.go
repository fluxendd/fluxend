package pkg

import (
	"github.com/rs/zerolog/log"
	"os/exec"
)

func ExecuteCommand(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Str("command", cmd.String()).
			Str("output", string(output)).
			Str("error", err.Error()).
			Msg("Command failed")

		return err
	}

	// Uncomment for debugging purposes
	//log.Info().
	//	Str("command", cmd.String()).
	//	Str("output", string(output)).
	//	Msg("Command successful")

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
