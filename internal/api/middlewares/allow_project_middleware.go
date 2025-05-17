package middlewares

import (
	"fluxton/internal/api/response"
	"fluxton/services"
	"github.com/labstack/echo/v4"
)

func AllowProjectMiddleware(settingService services.SettingService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool(c, "allowProjects") {
				return response.ForbiddenResponse(c, "project.error.disabled")
			}

			return next(c)
		}
	}
}
