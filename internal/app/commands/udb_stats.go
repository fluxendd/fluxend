package commands

import (
	"fluxend/internal/app"
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/stats"
	"fluxend/pkg"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

var udbStats = &cobra.Command{
	Use:   "udb.stats [database_name]",
	Short: "Pull stats from given database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectUUIDInput := args[0]
		projectUUID, err := uuid.Parse(projectUUIDInput)
		if err != nil {
			cmd.Printf("Invalid project UUID: %s\n", projectUUIDInput)
		}

		stats, err := getDatabaseStats(projectUUID)
		if err != nil {
			return err
		}

		if stats.DatabaseName == "" {
			cmd.Printf("Database %s not found", projectUUID)

			return nil
		}

		pkg.DumpJSON(stats)

		return nil
	},
}

func getDatabaseStats(projectUUID uuid.UUID) (stats.Stat, error) {
	container := app.InitializeContainer()

	authUser := auth.User{
		Uuid:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		RoleID: 1,
	}

	databaseStatsService := do.MustInvoke[stats.Service](container)

	pulledStats, err := databaseStatsService.GetAll(projectUUID, authUser)
	if err != nil {
		return stats.Stat{}, err
	}

	return pulledStats, nil
}
