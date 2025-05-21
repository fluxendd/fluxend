package middlewares

import (
	"fluxton/internal/api/response"
	"fluxton/internal/domain/setting"
	"github.com/labstack/echo/v4"
)

func AllowFormMiddleware(settingService setting.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool("allowForms") {
				return response.ForbiddenResponse(c, "form.error.disabled")
			}

			return next(c)
		}
	}
}
