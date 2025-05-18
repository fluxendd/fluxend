package commands

import (
	"fluxton/internal/app"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/database/stat"
	"fluxton/internal/domain/stats"
	"fluxton/pkg"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

var udbStats = &cobra.Command{
	Use:   "udb.stats [database_name]",
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

		pkg.DumpJSON(stats)

		return nil
	},
}

func getDatabaseStats(databaseName string) (stat.DatabaseStat, error) {
	container := app.InitializeContainer()

	authUser := auth.User{
		Uuid:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		RoleID: 1,
	}

	databaseStatsService := do.MustInvoke[stats.Service](container)

	pulledStats, err := databaseStatsService.GetAll(databaseName, authUser)
	if err != nil {
		return stat.DatabaseStat{}, err
	}

	return pulledStats, nil
}
