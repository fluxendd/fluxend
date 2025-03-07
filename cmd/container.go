package cmd

import (
	"fluxton/controllers"
	"fluxton/database"
	"fluxton/database/factories"
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
	do.Provide(injector, repositories.NewUserRepository)
	do.Provide(injector, repositories.NewRequestLogRepository)
	do.Provide(injector, repositories.NewDatabaseRepository)
	do.Provide(injector, repositories.NewSettingRepository)
	do.Provide(injector, repositories.NewCoreTableRepository)
	do.Provide(injector, repositories.NewOrganizationRepository)
	do.Provide(injector, repositories.NewProjectRepository)
	do.Provide(injector, repositories.NewFormRepository)
	do.Provide(injector, repositories.NewFormFieldRepository)
	do.Provide(injector, repositories.NewFormResponseRepository)
	do.Provide(injector, repositories.NewBucketRepository)
	do.Provide(injector, repositories.NewFileRepository)

	// Factories
	do.Provide(injector, factories.NewUserFactory)

	// policies
	do.Provide(injector, policies.NewOrganizationPolicy)
	do.Provide(injector, policies.NewProjectPolicy)

	// Services
	do.Provide(injector, services.NewUserService)
	do.Provide(injector, services.NewSettingService)
	do.Provide(injector, services.NewHealthService)
	do.Provide(injector, services.NewConnectionService)
	do.Provide(injector, services.NewOrganizationService)
	do.Provide(injector, services.NewProjectService)
	do.Provide(injector, services.NewTableService)
	do.Provide(injector, services.NewColumnService)
	do.Provide(injector, services.NewIndexService)
	do.Provide(injector, services.NewFunctionService)
	do.Provide(injector, services.NewFormService)
	do.Provide(injector, services.NewFormFieldService)
	do.Provide(injector, services.NewFormResponseService)
	do.Provide(injector, services.NewBucketService)
	do.Provide(injector, services.NewFileService)

	// Controllers
	do.Provide(injector, controllers.NewUserController)
	do.Provide(injector, controllers.NewSettingController)
	do.Provide(injector, controllers.NewHealthController)
	do.Provide(injector, controllers.NewOrganizationController)
	do.Provide(injector, controllers.NewOrganizationMemberController)
	do.Provide(injector, controllers.NewProjectController)
	do.Provide(injector, controllers.NewTableController)
	do.Provide(injector, controllers.NewColumnController)
	do.Provide(injector, controllers.NewIndexController)
	do.Provide(injector, controllers.NewFunctionController)
	do.Provide(injector, controllers.NewFormController)
	do.Provide(injector, controllers.NewFormFieldController)
	do.Provide(injector, controllers.NewFormResponseController)
	do.Provide(injector, controllers.NewBucketController)
	do.Provide(injector, controllers.NewFileController)

	return injector
}
