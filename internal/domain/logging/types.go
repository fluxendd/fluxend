package logging

import (
	"github.com/guregu/null/v6"
)

type ListInput struct {
	UserUuid  null.String `query:"userUuid"`
	Status    null.String `query:"status"`
	Method    null.String `query:"method"`
	Endpoint  null.String `query:"endpoint"`
	IPAddress null.String `query:"ipAddress"`
}
