package repositories

import (
	"fluxton/internal/domain/logging"
	"fluxton/internal/domain/shared"
	"fluxton/pkg"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type RequestLogRepository struct {
	db *sqlx.DB
}

func NewRequestLogRepository(injector *do.Injector) (logging.Repository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &RequestLogRepository{db: db}, nil
}

func (r *RequestLogRepository) List(paginationParams shared.PaginationParams) ([]logging.RequestLog, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `SELECT * FROM fluxton.api_logs ORDER BY :sort DESC LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

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
