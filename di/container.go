package di

import (
	"fluxton/controllers"
	"fluxton/db"
	"fluxton/factories"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/services"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

func InitializeContainer() *do.Injector {
	injector := do.New()

	// Database
	db.InitDB()
	do.Provide(injector, func(i *do.Injector) (*sqlx.DB, error) {
		return db.GetDB(), nil
	})

	// Repositories
	do.Provide(injector, repositories.NewUserRepository)
	do.Provide(injector, repositories.NewDatabaseRepository)
	do.Provide(injector, repositories.NewSettingRepository)
	do.Provide(injector, repositories.NewCoreTableRepository)
	do.Provide(injector, repositories.NewOrganizationRepository)
	do.Provide(injector, repositories.NewProjectRepository)

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
	do.Provide(injector, services.NewRowService)
	do.Provide(injector, services.NewIndexService)

	// Controllers
	do.Provide(injector, controllers.NewUserController)
	do.Provide(injector, controllers.NewSettingController)
	do.Provide(injector, controllers.NewHealthController)
	do.Provide(injector, controllers.NewOrganizationController)
	do.Provide(injector, controllers.NewOrganizationUserController)
	do.Provide(injector, controllers.NewProjectController)
	do.Provide(injector, controllers.NewTableController)
	do.Provide(injector, controllers.NewColumnController)
	do.Provide(injector, controllers.NewRowController)
	do.Provide(injector, controllers.NewIndexController)

	return injector
}
