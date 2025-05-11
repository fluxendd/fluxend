package email

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type SESService interface {
	SendEmail(to, subject, body string) error
}

type SESServiceImpl struct {
	client *ses.Client
}

func NewSESService() (SESService, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := ses.NewFromConfig(cfg)

	return &SESServiceImpl{
		client: client,
	}, nil
}

func (s *SESServiceImpl) SendEmail(to, subject, body string) error {
	input := &ses.SendEmailInput{
		Source: aws.String(os.Getenv("SES_EMAIL_SOURCE")),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(body),
				},
			},
		},
	}

	_, err := s.client.SendEmail(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
