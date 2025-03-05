package middlewares

import (
	"fluxton/models"
	"fluxton/repositories"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQL driver
	"time"
)

func RequestLoggerMiddleware(requestLogRepo *repositories.RequestLogRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := next(c)
			authUserUUID, _ := utils.NewAuth(c).Uuid()
			APIKey, err := utils.GetUUIDQueryParam(c, "api_key", true)
			if err != nil {
				APIKey = uuid.MustParse("00000000-0000-0000-0000-000000000000")
			}

			request := c.Request()

			logEntry := models.RequestLog{
				UserUuid:  authUserUUID,
				APIKey:    APIKey, // TODO: implement API key parsing/reading later
				Method:    request.Method,
				Endpoint:  request.URL.Path,
				IPAddress: c.RealIP(),
				UserAgent: request.UserAgent(),
				Params:    request.URL.RawQuery,
				Body:      "{}", // TODO: implement body parsing/reading later
				CreatedAt: time.Now(),
			}

			go requestLogRepo.Create(&logEntry)

			return res
		}
	}
}
