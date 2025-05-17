package routes

import (
	"fluxton/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterBackup(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc, allowBackupMiddleware echo.MiddlewareFunc) {
	backupController := do.MustInvoke[*handlers.BackupHandler](container)

	formsGroup := e.Group("api/backups", authMiddleware, allowBackupMiddleware)

	formsGroup.POST("", backupController.Store)
	formsGroup.GET("", backupController.List)
	formsGroup.GET("/:backupUUID", backupController.Show)
	formsGroup.DELETE("/:backupUUID", backupController.Delete)
}
