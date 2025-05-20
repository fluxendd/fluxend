package email

import (
	"fluxton/internal/config/constants"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type Provider interface {
	Send(ctx echo.Context, to, subject, body string) error
}

type Factory struct {
	injector *do.Injector
}

func NewFactory(injector *do.Injector) (*Factory, error) {
	return &Factory{injector: injector}, nil
}

func (f *Factory) CreateProvider(ctx echo.Context, providerType string) (Provider, error) {
	switch providerType {
	case constants.EmailDriverSES:
		return NewSESProvider(ctx, f.injector)
	case constants.EmailDriverSendGrid:
		return NewSendGridProvider(ctx, f.injector)
	case constants.EmailDriverMailgun:
		return NewMailgunProvider(ctx, f.injector)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", providerType)
	}
}
