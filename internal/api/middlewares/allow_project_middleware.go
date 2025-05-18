package middlewares

import (
	"fluxton/internal/api/response"
	"fluxton/internal/domain/setting"
	"github.com/labstack/echo/v4"
)

func AllowProjectMiddleware(settingService setting.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool(c, "allowProjects") {
				return response.ForbiddenResponse(c, "project.error.disabled")
			}

			return next(c)
		}
	}
}
