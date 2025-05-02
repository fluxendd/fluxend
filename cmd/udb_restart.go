package cmd

import (
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/services"
	"fmt"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

var udbRestart = &cobra.Command{
	Use:   "udb.restart",
	Short: "Restart all PostGREST instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		return restartPostgrestInstances()
	},
}

func restartPostgrestInstances() error {
	container := InitializeContainer()

	// Inject dependencies
	projectRepository := do.MustInvoke[*repositories.ProjectRepository](container)
	postgrestService := do.MustInvoke[services.PostgrestService](container)

	projects, err := projectRepository.List(requests.PaginationParams{Page: 1, Limit: 1000})
	if err != nil {
		return fmt.Errorf("error fetching projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("No projects found")

		return nil
	}

	fmt.Printf("Found %d projects\n", len(projects))

	// Restart PostGREST instances for each project
	for i, project := range projects {
		if project.DBName == "" {
			continue
		}

		// Log restart process for each project
		fmt.Printf("Restarting PostGREST instance for project %s (%d/%d)\n", project.DBName, i+1, len(projects))

		hasContainer := postgrestService.HasContainer(project.DBName)
		if hasContainer {
			postgrestService.RemoveContainer(project.DBName)
		}

		postgrestService.StartContainer(project.DBName)
	}

	return nil
}
