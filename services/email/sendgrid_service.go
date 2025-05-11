package email

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridServiceImpl struct {
	client *sendgrid.Client
}

func NewSendGridService() (EmailInterface, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SENDGRID_API_KEY environment variable is required")
	}

	client := sendgrid.NewSendClient(apiKey)

	return &SendGridServiceImpl{
		client: client,
	}, nil
}

func (s *SendGridServiceImpl) Send(to, subject, body string) error {
	from := os.Getenv("SENDGRID_EMAIL_SOURCE")
	if from == "" {
		return fmt.Errorf("SENDGRID_EMAIL_SOURCE environment variable is required")
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
