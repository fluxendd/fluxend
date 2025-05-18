package app

import (
	"fluxton/internal/adapters/client"
	"fluxton/internal/adapters/postgrest"
	"fluxton/internal/api/handlers"
	"fluxton/internal/database"
	"fluxton/internal/database/repositories"
	"fluxton/internal/domain/backup"
	"fluxton/internal/domain/database/column"
	"fluxton/internal/domain/database/function"
	"fluxton/internal/domain/database/index"
	"fluxton/internal/domain/database/table"
	"fluxton/internal/domain/file_import"
	"fluxton/internal/domain/form"
	"fluxton/internal/domain/health"
	"fluxton/internal/domain/organization"
	"fluxton/internal/domain/project"
	"fluxton/internal/domain/setting"
	"fluxton/internal/domain/stats"
	"fluxton/internal/domain/storage/container"
	"fluxton/internal/domain/storage/file"
	"fluxton/internal/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

func InitializeContainer() *do.Injector {
	injector := do.New()

	// Database
	database.InitDB()
	do.Provide(injector, func(i *do.Injector) (*sqlx.DB, error) {
		return database.GetDB(), nil
	})

	// Repositories
	do.Provide(injector, repositories.NewUserRepository)
	do.Provide(injector, repositories.NewRequestLogRepository)
	do.Provide(injector, client.NewClientRepository)
	do.Provide(injector, repositories.NewSettingRepository)
	do.Provide(injector, repositories.NewOrganizationRepository)
	do.Provide(injector, repositories.NewProjectRepository)
	do.Provide(injector, repositories.NewFormRepository)
	do.Provide(injector, repositories.NewFormFieldRepository)
	do.Provide(injector, repositories.NewFormResponseRepository)
	do.Provide(injector, repositories.NewContainerRepository)
	do.Provide(injector, repositories.NewFileRepository)
	do.Provide(injector, repositories.NewBackupRepository)

	// Factories
	//do.Provide(injector, factories.NewUserFactory)

	// policies
	do.Provide(injector, organization.NewOrganizationPolicy)
	do.Provide(injector, project.NewProjectPolicy)

	// Services
	do.Provide(injector, user.NewUserService)
	do.Provide(injector, setting.NewSettingService)
	do.Provide(injector, health.NewHealthService)
	do.Provide(injector, client.NewClientService)
	do.Provide(injector, stats.NewDatabaseStatsService)
	do.Provide(injector, organization.NewOrganizationService)
	do.Provide(injector, project.NewProjectService)
	do.Provide(injector, postgrest.NewPostgrestService)
	do.Provide(injector, table.NewTableService)
	do.Provide(injector, file_import.NewFileImportService)
	do.Provide(injector, column.NewColumnService)
	do.Provide(injector, index.NewIndexService)
	do.Provide(injector, function.NewFunctionService)
	do.Provide(injector, form.NewFormService)
	do.Provide(injector, form.NewFormFieldValidationService)
	do.Provide(injector, form.NewFieldService)
	do.Provide(injector, form.NewFormResponseService)
	do.Provide(injector, container.NewContainerService)
	do.Provide(injector, file.NewFileService)
	do.Provide(injector, backup.NewBackupWorkflowService)
	do.Provide(injector, backup.NewBackupService)

	// Handlers
	do.Provide(injector, handlers.NewUserHandler)
	do.Provide(injector, handlers.NewSettingHandler)
	do.Provide(injector, handlers.NewHealthHandler)
	do.Provide(injector, handlers.NewOrganizationHandler)
	do.Provide(injector, handlers.NewOrganizationMemberHandler)
	do.Provide(injector, handlers.NewProjectHandler)
	do.Provide(injector, handlers.NewTableHandler)
	do.Provide(injector, handlers.NewColumnHandler)
	do.Provide(injector, handlers.NewIndexHandler)
	do.Provide(injector, handlers.NewFunctionHandler)
	do.Provide(injector, handlers.NewFormHandler)
	do.Provide(injector, handlers.NewFormFieldHandler)
	do.Provide(injector, handlers.NewFormResponseHandler)
	do.Provide(injector, handlers.NewContainerHandler)
	do.Provide(injector, handlers.NewFileHandler)
	do.Provide(injector, handlers.NewBackupHandler)

	return injector
}
