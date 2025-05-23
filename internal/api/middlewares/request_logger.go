package middlewares

import (
	"bytes"
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/logging"
	"fluxton/pkg/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

func RequestLogger(requestLogRepo logging.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			requestBody := readBody(request)

			res := next(c)
			authUserUUID, _ := auth.NewAuth(c).Uuid()

			logEntry := logging.RequestLog{
				UserUuid:  authUserUUID,
				APIKey:    uuid.MustParse("00000000-0000-0000-0000-000000000000"), // TODO: implement API key parsing/reading later
				Method:    request.Method,
				Status:    c.Response().Status,
				Endpoint:  request.URL.Path,
				IPAddress: c.RealIP(),
				UserAgent: request.UserAgent(),
				Params:    request.URL.RawQuery,
				Body:      requestBody,
				CreatedAt: time.Now(),
			}

			log.Info().
				Str("action", constants.ActionAPIRequest).
				Str("user_uuid", authUserUUID.String()).
				Str("method", logEntry.Method).
				Str("endpoint", logEntry.Endpoint).
				Str("ip_address", logEntry.IPAddress).
				Str("user_agent", logEntry.UserAgent).
				Int("status", c.Response().Status).
				Msg("")

			go requestLogRepo.Create(&logEntry)

			return res
		}
	}
}

func readBody(r *http.Request) string {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")

		return ""
	}

	// Restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}
