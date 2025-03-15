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
	Run:   runMakeModelCmd,
}

func runMakeModelCmd(cmd *cobra.Command, args []string) {
	modelName := args[0]

	if !isValidModelName(modelName) {
		fmt.Println("Error: Model name can only contain letters and underscores.")
		return
	}

	if err := createModelFiles(modelName); err != nil {
		fmt.Println("Error:", err)
	}
}

func isValidModelName(modelName string) bool {
	return validModelName.MatchString(modelName)
}

func createModelFiles(modelName string) error {
	modelStubPath := "stubs/model.stub"
	modelsDir := "models"
	modelOutputFile := filepath.Join(modelsDir, strings.ToLower(modelName)+".go")

	resourceStubPath := "stubs/resource.stub"
	resourcesDir := "resources"
	resourceOutputFile := filepath.Join(resourcesDir, strings.ToLower(modelName)+"Resource.go")

	if err := ensureDirExists(modelsDir); err != nil {
		return fmt.Errorf("creating models directory: %w", err)
	}

	if err := ensureFileDoesNotExist(modelOutputFile); err != nil {
		return fmt.Errorf("checking if model exists: %w", err)
	}

	if err := writeFileWithReplacedContent(modelStubPath, modelOutputFile, modelName); err != nil {
		return fmt.Errorf("writing model file: %w", err)
	}

	if err := writeFileWithReplacedContent(resourceStubPath, resourceOutputFile, modelName); err != nil {
		return fmt.Errorf("writing resource file: %w", err)
	}

	fmt.Printf("Model %s created successfully at %s\n", modelName, modelOutputFile)
	return nil
}

func ensureDirExists(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func ensureFileDoesNotExist(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists at %s", filePath)
	}
	return nil
}

func writeFileWithReplacedContent(stubPath, outputFile, modelName string) error {
	stubContent, err := os.ReadFile(stubPath)
	if err != nil {
		return fmt.Errorf("reading stub file %s: %w", stubPath, err)
	}

	replacedContent := replacePlaceholders(string(stubContent), modelName)

	if err := os.WriteFile(outputFile, []byte(replacedContent), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", outputFile, err)
	}

	return nil
}

func replacePlaceholders(content, modelName string) string {
	content = strings.ReplaceAll(content, "{{modelName}}", modelName)
	content = strings.ReplaceAll(content, "{{modelLowercase}}", strings.ToLower(modelName))
	return content
}
