package logging

import (
	"github.com/google/uuid"
)

type Response struct {
	Uuid      uuid.UUID `json:"uuid"`
	UserUuid  uuid.UUID `json:"userUuid"`
	Method    string    `json:"method"`
	Status    int       `json:"status"`
	Endpoint  string    `json:"endpoint"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	Params    string    `json:"params"`
	Body      string    `json:"body"`
	CreatedAt string    `json:"createdAt"`
}
