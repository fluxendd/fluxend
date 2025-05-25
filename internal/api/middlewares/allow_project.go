package middlewares

import (
	"fluxend/internal/api/response"
	"fluxend/internal/domain/setting"
	"github.com/labstack/echo/v4"
)

func AllowProject(settingService setting.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool("allowProjects") {
				return response.ForbiddenResponse(c, "project.error.disabled")
			}

			return next(c)
		}
	}
}
