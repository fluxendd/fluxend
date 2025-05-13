package email

import (
	"fluxton/constants"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func GetProvider(ctx echo.Context, injector *do.Injector, provider string) (EmailInterface, error) {
	switch provider {
	case constants.EmailDriverSES:
		return NewSESService(ctx, injector)
	case constants.EmailDriverSendGrid:
		return NewSendGridService(ctx, injector)
	case constants.EmailDriverMailgun:
		return NewMailgunService(ctx, injector)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", provider)
	}
}
