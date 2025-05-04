package middlewares

import (
	"fluxton/responses"
	"fluxton/services"
	"github.com/labstack/echo/v4"
)

func AllowBackupMiddleware(settingService services.SettingService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool("allowBackup") {
				return responses.ForbiddenResponse(c, "backup.error.disabled")
			}

			return next(c)
		}
	}
}
