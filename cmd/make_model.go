package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var validModelName = regexp.MustCompile(`^[A-Za-z_]+$`)

var makeModelCmd = &cobra.Command{
	Use:   "make:model [model_name]",
	Short: "Creates a new model file in the models directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]

		if !validModelName.MatchString(modelName) {
			fmt.Println("Error: Model name can only contain letters and underscores.")

			return
		}

		makeModel(modelName)
	},
}

func makeModel(modelName string) {
	stubPath := "stubs/model.stub"
	modelsDir := "models"
	outputFile := filepath.Join(modelsDir, strings.ToLower(modelName)+".go")

	if err := os.MkdirAll(modelsDir, os.ModePerm); err != nil {
		fmt.Println("Error creating models directory:", err)
		return
	}

	if _, err := os.Stat(outputFile); err == nil {
		fmt.Printf("Error: Model %s already exists at %s\n", modelName, outputFile)
		return
	}

	stubContent, err := os.ReadFile(stubPath)
	if err != nil {
		fmt.Println("Error reading stub file:", err)
		return
	}

	replacedContent := strings.ReplaceAll(string(stubContent), "{{modelName}}", modelName)

	if err := os.WriteFile(outputFile, []byte(replacedContent), 0644); err != nil {
		fmt.Println("Error writing model file:", err)
		return
	}

	fmt.Printf("Model %s created successfully at %s\n", modelName, outputFile)
}
