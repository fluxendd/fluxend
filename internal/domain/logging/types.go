package logging

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"time"
)

type ListInput struct {
	ProjectUuid uuid.NullUUID `query:"projectUuid"`
	UserUuid    uuid.NullUUID `query:"userUuid"`
	Status      null.String   `query:"status"`
	Method      null.String   `query:"method"`
	Endpoint    null.String   `query:"endpoint"`
	IPAddress   null.String   `query:"ipAddress"`
	StartTime   time.Time     `query:"startTime"`
	EndTime     time.Time     `query:"endTime"`
}

type StoreInput struct {
	Endpoint    string    `json:"endpoint"`
	DbName      string    `json:"DbName"`
	ProjectUUID uuid.UUID `json:"projectUuid"`
	IPAddress   string    `json:"ipAddress"`
	Host        string    `json:"host"`
	Method      string    `json:"method"`
}
