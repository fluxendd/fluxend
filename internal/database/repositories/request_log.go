package repositories

import (
	"fluxend/internal/domain/logging"
	"fluxend/internal/domain/shared"
	"fmt"
	"github.com/samber/do"
	"strings"
)

type RequestLogRepository struct {
	db shared.DB
}

func NewRequestLogRepository(injector *do.Injector) (logging.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &RequestLogRepository{db: db}, nil
}

func (r *RequestLogRepository) List(input *logging.ListInput, paginationParams shared.PaginationParams) ([]logging.RequestLog, error) {
	var filters []string

	offset := (paginationParams.Page - 1) * paginationParams.Limit
	params := map[string]interface{}{
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	if input.ProjectUuid.Valid {
		filters = append(filters, "project_uuid = :project_uuid")
		params["project_uuid"] = input.ProjectUuid.UUID.String()
	}

	if input.UserUuid.Valid {
		filters = append(filters, "user_uuid = :user_uuid")
		params["user_uuid"] = input.UserUuid.UUID.String()
	}

	if input.Status.Valid {
		filters = append(filters, "status = :status")
		params["status"] = input.Status
	}

	if input.Method.Valid {
		filters = append(filters, "method = :method")
		params["method"] = input.Method
	}

	if input.Endpoint.Valid {
		filters = append(filters, "endpoint = :endpoint")
		params["endpoint"] = input.Endpoint
	}

	if input.IPAddress.Valid {
		filters = append(filters, "ip_address = :ip_address")
		params["ip_address"] = input.IPAddress
	}

	query := `SELECT * FROM fluxend.api_logs`
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"status":     true,
		"method":     true,
	}
	sortColumn := "created_at" // default
	if allowedSorts[paginationParams.Sort] {
		sortColumn = paginationParams.Sort
	}

	query += fmt.Sprintf(" ORDER BY %s DESC LIMIT :limit OFFSET :offset", sortColumn)

	var requestLogs []logging.RequestLog
	return requestLogs, r.db.SelectNamedList(&requestLogs, query, params)
}

func (r *RequestLogRepository) Create(requestLog *logging.RequestLog) (*logging.RequestLog, error) {
	return requestLog, r.db.WithTransaction(func(tx shared.Tx) error {
		query := `
        INSERT INTO fluxend.api_logs (
            user_uuid, api_key, method, status, endpoint, ip_address, user_agent, params, body
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
        RETURNING uuid
        `

		return tx.QueryRowx(
			query,
			requestLog.UserUuid,
			requestLog.APIKey,
			requestLog.Method,
			requestLog.Status,
			requestLog.Endpoint,
			requestLog.IPAddress,
			requestLog.UserAgent,
			requestLog.Params,
			requestLog.Body,
		).Scan(&requestLog.Uuid)
	})
}
