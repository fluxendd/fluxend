package middlewares

import (
	"fluxton/models"
	"fluxton/repositories"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/zerolog/log"
	"io"
	"time"
)

func RequestLoggerMiddleware(requestLogRepo *repositories.RequestLogRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := next(c)
			authUserUUID, _ := utils.NewAuth(c).Uuid()

			request := c.Request()

			logEntry := models.RequestLog{
				UserUuid:  authUserUUID,
				APIKey:    uuid.MustParse("00000000-0000-0000-0000-000000000000"), // TODO: implement API key parsing/reading later
				Method:    request.Method,
				Endpoint:  request.URL.Path,
				IPAddress: c.RealIP(),
				UserAgent: request.UserAgent(),
				Params:    request.URL.RawQuery,
				Body:      readBody(request.Body), // TODO: look into streams and buffering
				CreatedAt: time.Now(),
			}

			log.Log().
				Str("user_uuid", authUserUUID.String()).
				Str("method", logEntry.Method).
				Str("endpoint", logEntry.Endpoint).
				Str("ip_address", logEntry.IPAddress).
				Str("user_agent", logEntry.UserAgent).
				Msg("incoming request")

			go requestLogRepo.Create(&logEntry)

			return res
		}
	}
}

func readBody(body io.ReadCloser) string {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Printf("failed to read request body: %v", err)

		return ""
	}

	return string(data)
}
