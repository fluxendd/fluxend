package email

import (
	"fluxton/internal/config/constants"
	"fmt"
	"github.com/samber/do"
)

type Provider interface {
	Send(to, subject, body string) error
}

type Factory struct {
	injector *do.Injector
}

func NewFactory(injector *do.Injector) (*Factory, error) {
	return &Factory{injector: injector}, nil
}

func (f *Factory) CreateProvider(providerType string) (Provider, error) {
	switch providerType {
	case constants.EmailDriverSES:
		return NewSESProvider(f.injector)
	case constants.EmailDriverSendGrid:
		return NewSendGridProvider(f.injector)
	case constants.EmailDriverMailgun:
		return NewMailgunProvider(f.injector)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", providerType)
	}
}
