package middlewares

import (
	"fluxton/responses"
	"fluxton/services"
	"github.com/labstack/echo/v4"
)

func FormEnabledMiddleware(settingService services.SettingService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := settingService.ValidateFormsEnabled(); err != nil {
				return responses.ForbiddenResponse(c, "form.error.disabled")
			}

			return next(c)
		}
	}
}
