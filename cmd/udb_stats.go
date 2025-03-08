package cmd

import (
	"fluxton/models"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

var udbStats = &cobra.Command{
	Use:   "udb:stats [database_name]",
	Short: "Pull stats from given database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		databaseName := args[0]

		stats, err := getDatabaseStats(databaseName)
		if err != nil {
			return err
		}

		if stats.DatabaseName == "" {
			cmd.Printf("Database %s not found", databaseName)

			return nil
		}

		utils.DumpJSON(stats)

		return nil
	},
}

func getDatabaseStats(databaseName string) (models.DatabaseStat, error) {
	container := InitializeContainer()

	authUser := models.AuthUser{
		Uuid:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		RoleID: 1,
	}

	databaseStatsService := do.MustInvoke[services.DatabaseStatsService](container)

	stats, err := databaseStatsService.GetAll(databaseName, authUser)
	if err != nil {
		return models.DatabaseStat{}, err
	}

	return stats, nil
}
