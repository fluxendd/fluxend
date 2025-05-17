package app

import (
	"fluxton/controllers"
	"fluxton/internal/adapters/connection"
	"fluxton/internal/adapters/postgrest"
	"fluxton/internal/api/handlers"
	"fluxton/internal/database"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/organization"
	"fluxton/internal/domain/user"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/services"
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
	do.Provide(injector, repositories2.NewUserRepository)
	do.Provide(injector, repositories.NewRequestLogRepository)
	do.Provide(injector, repositories.NewDatabaseRepository)
	do.Provide(injector, repositories.NewSettingRepository)
	do.Provide(injector, repositories2.NewOrganizationRepository)
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
	do.Provide(injector, policies.NewProjectPolicy)

	// Services
	do.Provide(injector, user.NewUserService)
	do.Provide(injector, services.NewSettingService)
	do.Provide(injector, services.NewHealthService)
	do.Provide(injector, connection.NewConnectionService)
	do.Provide(injector, services.NewDatabaseStatsService)
	do.Provide(injector, organization.NewOrganizationService)
	do.Provide(injector, services.NewProjectService)
	do.Provide(injector, postgrest.NewPostgrestService)
	do.Provide(injector, services.NewTableService)
	do.Provide(injector, services.NewFileImportService)
	do.Provide(injector, services.NewColumnService)
	do.Provide(injector, services.NewIndexService)
	do.Provide(injector, services.NewFunctionService)
	do.Provide(injector, services.NewFormService)
	do.Provide(injector, services.NewFormFieldValidationService)
	do.Provide(injector, services.NewFormFieldService)
	do.Provide(injector, services.NewFormResponseService)
	do.Provide(injector, services.NewContainerService)
	do.Provide(injector, services.NewFileService)
	do.Provide(injector, services.NewBackupWorkflowService)
	do.Provide(injector, services.NewBackupService)

	// Handlers
	do.Provide(injector, handlers.NewUserHandler)
	do.Provide(injector, controllers.NewSettingController)
	do.Provide(injector, controllers.NewHealthController)
	do.Provide(injector, handlers.NewOrganizationHandler)
	do.Provide(injector, handlers.NewOrganizationMemberHandler)
	do.Provide(injector, controllers.NewProjectController)
	do.Provide(injector, controllers.NewTableController)
	do.Provide(injector, controllers.NewColumnController)
	do.Provide(injector, controllers.NewIndexController)
	do.Provide(injector, controllers.NewFunctionController)
	do.Provide(injector, controllers.NewFormController)
	do.Provide(injector, controllers.NewFormFieldController)
	do.Provide(injector, controllers.NewFormResponseController)
	do.Provide(injector, controllers.NewContainerController)
	do.Provide(injector, controllers.NewFileController)
	do.Provide(injector, controllers.NewBackupController)

	return injector
}
