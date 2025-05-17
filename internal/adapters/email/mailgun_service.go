package email

import (
	"context"
	"fluxton/internal/domain/setting"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunServiceImpl struct {
	client         *mailgun.MailgunImpl
	domain         string
	settingService setting.SettingService
}

func NewMailgunService(ctx echo.Context, injector *do.Injector) (EmailInterface, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	apiKey := settingService.GetValue(ctx, "mailgunApiKey")
	if apiKey == "" {
		return nil, fmt.Errorf("mailgun API key is required")
	}

	domain := settingService.GetValue(ctx, "mailgunDomain")
	if domain == "" {
		return nil, fmt.Errorf("mailgun domain is required")
	}

	userSelectedRegion := settingService.Get(ctx, "mailgunRegion")
	mailgunRegion := userSelectedRegion.Value
	if mailgunRegion == "" {
		mailgunRegion = userSelectedRegion.DefaultValue
	}

	client := mailgun.NewMailgun(domain, apiKey)

	if mailgunRegion == "eu" {
		client.SetAPIBase(mailgun.APIBaseEU)
	}

	return &MailgunServiceImpl{
		client:         client,
		domain:         domain,
		settingService: settingService,
	}, nil
}

func (m *MailgunServiceImpl) Send(ctx echo.Context, to, subject, body string) error {
	from := m.settingService.GetValue(ctx, "mailgunEmailSource")
	if from == "" {
		return fmt.Errorf("mailgunEmailSource is required")
	}

	message := m.client.NewMessage(
		from,
		subject,
		body,
		to,
	)

	backgroundContext, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, _, err := m.client.Send(backgroundContext, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	if resp != "" {
		// Mailgun returns an empty string on success
		return fmt.Errorf("unexpected response: %s", resp)
	}

	return nil
}
