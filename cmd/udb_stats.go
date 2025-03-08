package cmd

import (
	"fluxton/models"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

// udbStats runs in background and populates table with stats
var udbStats = &cobra.Command{
	Use:   "udb:stats [project_uuid]",
	Short: "Pull stats from give project's database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectUUIDString := args[0]
		getDatabaseStats(projectUUIDString)
	},
}

func getDatabaseStats(projectUUIDString string) {
	container := InitializeContainer()

	authUser := models.AuthUser{
		Uuid:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		RoleID: 1,
	}

	projectUUID, err := uuid.Parse(projectUUIDString)
	if err != nil {
		utils.DumpJSON(err)
	}

	databaseStatsService := do.MustInvoke[services.DatabaseStatsService](container)
	projectService := do.MustInvoke[services.ProjectService](container)

	databaseName, err := projectService.GetDatabaseNameByUUID(projectUUID, authUser)
	if err != nil {
		utils.DumpJSON(err)
	}

	stats, err := databaseStatsService.GetSizePerTable(databaseName, authUser)
	if err != nil {
		utils.DumpJSON(err)
	}

	utils.DumpJSON(stats)
}
