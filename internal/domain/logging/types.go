package logging

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type ListInput struct {
	ProjectUuid uuid.NullUUID `query:"projectUuid"`
	UserUuid    uuid.NullUUID `query:"userUuid"`
	Status      null.String   `query:"status"`
	Method      null.String   `query:"method"`
	Endpoint    null.String   `query:"endpoint"`
	IPAddress   null.String   `query:"ipAddress"`
}
