package middlewares

import (
	"fluxton/internal/api/response"
	"fluxton/services"
	"github.com/labstack/echo/v4"
)

func AllowStorageMiddleware(settingService services.SettingService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool(c, "allowStorage") {
				return response.ForbiddenResponse(c, "storage.error.disabled")
			}

			return next(c)
		}
	}
}
