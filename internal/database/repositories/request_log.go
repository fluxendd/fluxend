package repositories

import (
	"fluxton/internal/domain/logging"
	"fluxton/internal/domain/shared"
	"fluxton/pkg"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"strings"
)

type RequestLogRepository struct {
	db *sqlx.DB
}

func NewRequestLogRepository(injector *do.Injector) (logging.Repository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &RequestLogRepository{db: db}, nil
}

func (r *RequestLogRepository) List(input *logging.ListInput, paginationParams shared.PaginationParams) ([]logging.RequestLog, error) {
	var filters []string

	offset := (paginationParams.Page - 1) * paginationParams.Limit
	params := map[string]interface{}{
		"limit":  paginationParams.Limit,
		"offset": offset,
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

	query := `SELECT * FROM fluxton.api_logs`
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

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var requestLogs []logging.RequestLog
	for rows.Next() {
		var form logging.RequestLog
		if err := rows.StructScan(&form); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		requestLogs = append(requestLogs, form)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return requestLogs, nil
}

func (r *RequestLogRepository) Create(requestLog *logging.RequestLog) (*logging.RequestLog, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	query := `
    INSERT INTO fluxton.api_logs (
        user_uuid, api_key, method, status, endpoint, ip_address, user_agent, params, body
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
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

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, pkg.FormatError(queryErr, "insert", pkg.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return requestLog, nil
}
