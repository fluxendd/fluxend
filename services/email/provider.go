package email

import (
	"fluxton/constants"
	"fmt"
)

func GetProvider(provider string) (EmailInterface, error) {
	switch provider {
	case constants.EmailDriverSES:
		return NewSESService()
	case constants.EmailDriverSendGrid:
		return NewSendGridService()
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", provider)
	}
}
