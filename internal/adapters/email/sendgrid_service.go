package email

import (
	"fluxton/internal/domain/setting"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridServiceImpl struct {
	client         *sendgrid.Client
	settingService setting.SettingService
}

func NewSendGridService(ctx echo.Context, injector *do.Injector) (EmailInterface, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	apiKey := settingService.GetValue(ctx, "sendgridApiKey")
	if apiKey == "" {
		return nil, fmt.Errorf("sendgridApiKey is required")
	}

	client := sendgrid.NewSendClient(apiKey)

	return &SendGridServiceImpl{
		client:         client,
		settingService: settingService,
	}, nil
}

func (s *SendGridServiceImpl) Send(ctx echo.Context, to, subject, body string) error {
	from := s.settingService.GetValue(ctx, "sendgridEmailSource")
	if from == "" {
		return fmt.Errorf("sendgridEmailSource is required")
	}

	message := mail.NewSingleEmail(
		mail.NewEmail("", from),
		subject,
		mail.NewEmail("", to),
		body,
		"", // HTML content is empty as we're using plain text
	)

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: status code %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
