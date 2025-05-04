package middlewares

import (
	"fluxton/responses"
	"fluxton/services"
	"github.com/labstack/echo/v4"
)

func AllowFormMiddleware(settingService services.SettingService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !settingService.GetBool("allowForm") {
				return responses.ForbiddenResponse(c, "form.error.disabled")
			}

			return next(c)
		}
	}
}
