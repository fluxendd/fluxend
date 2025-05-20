package email

import (
	"context"
	"fluxton/internal/domain/setting"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type SESServiceImpl struct {
	client         *ses.Client
	settingService setting.Service
}

func NewSESProvider(ctx echo.Context, injector *do.Injector) (Provider, error) {
	settingService, err := setting.NewSettingService(injector)
	if err != nil {
		return nil, err
	}

	accessKey := settingService.GetValue(ctx, "awsAccessKey")
	secretKey := settingService.GetValue(ctx, "awsSecretKey")
	region := settingService.GetValue(ctx, "awsRegion")

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
		client:         client,
		settingService: settingService,
	}, nil
}

func (s *SESServiceImpl) Send(ctx echo.Context, to, subject, body string) error {
	from := s.settingService.GetValue(ctx, "sesEmailSource")
	if from == "" {
		return fmt.Errorf("sesEmailSource is required")
	}

	input := &ses.SendEmailInput{
		Source: aws.String(from),
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
