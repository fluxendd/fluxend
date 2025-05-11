package email

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunServiceImpl struct {
	client *mailgun.MailgunImpl
	domain string
}

func NewMailgunService() (EmailInterface, error) {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MAILGUN_API_KEY environment variable is required")
	}

	domain := os.Getenv("MAILGUN_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("MAILGUN_DOMAIN environment variable is required")
	}

	// Optional region, defaults to US
	region := os.Getenv("MAILGUN_REGION")
	if region == "" {
		region = "us" // Default to US region
	}

	client := mailgun.NewMailgun(domain, apiKey)

	// Set region if specified
	if region == "eu" {
		client.SetAPIBase(mailgun.APIBaseEU)
	}

	return &MailgunServiceImpl{
		client: client,
		domain: domain,
	}, nil
}

func (m *MailgunServiceImpl) Send(to, subject, body string) error {
	from := os.Getenv("MAILGUN_EMAIL_SOURCE")
	if from == "" {
		return fmt.Errorf("MAILGUN_EMAIL_SOURCE environment variable is required")
	}

	message := m.client.NewMessage(
		from,
		subject,
		body,
		to,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := m.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	if resp != "" {
		// Mailgun returns an empty string on success
		return fmt.Errorf("unexpected response: %s", resp)
	}

	// Log the message ID if you want
	_ = id

	return nil
}
