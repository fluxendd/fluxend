package logging

import (
	"github.com/google/uuid"
	"time"
)

type RequestLog struct {
	Uuid      uuid.UUID `db:"uuid" json:"uuid"`
	UserUuid  uuid.UUID `db:"user_uuid" json:"userUuid"`
	APIKey    uuid.UUID `db:"api_key" json:"apiKey"`
	Method    string    `db:"method" json:"method"`
	Status    int       `db:"status" json:"status"`
	Endpoint  string    `db:"endpoint" json:"endpoint"`
	IPAddress string    `db:"ip_address" json:"ipAddress"`
	UserAgent string    `db:"user_agent" json:"userAgent"`
	Params    string    `db:"params" json:"params"`
	Body      string    `db:"body" json:"body"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
