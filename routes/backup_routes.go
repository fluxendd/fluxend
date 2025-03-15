package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterBackup(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	backupController := do.MustInvoke[*controllers.BackupController](container)

	formsGroup := e.Group("api/backups", authMiddleware)

	formsGroup.POST("", backupController.Store)
	formsGroup.GET("", backupController.List)
	formsGroup.GET("/:backupUUID", backupController.Show)
	formsGroup.DELETE("/:backupUUID", backupController.Delete)
}
