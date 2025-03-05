package repositories

import (
	"fluxton/models"
	"fluxton/utils"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type RequestLogRepository struct {
	db *sqlx.DB
}

func NewRequestLogRepository(injector *do.Injector) (*RequestLogRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &RequestLogRepository{db: db}, nil
}

func (r *RequestLogRepository) List(paginationParams utils.PaginationParams) ([]models.RequestLog, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `SELECT * FROM fluxton.request_logs ORDER BY :sort DESC LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var requestLogs []models.RequestLog
	for rows.Next() {
		var form models.RequestLog
		if err := rows.StructScan(&form); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		requestLogs = append(requestLogs, form)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return requestLogs, nil
}

func (r *RequestLogRepository) Create(requestLog *models.RequestLog) (*models.RequestLog, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `
    INSERT INTO fluxton.request_logs (
        user_uuid, api_key, method, endpoint, ip_address, user_agent, params
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		requestLog.UserUuid,
		requestLog.APIKey,
		requestLog.Method,
		requestLog.Endpoint,
		requestLog.IPAddress,
		requestLog.UserAgent,
		requestLog.Params,
	).Scan(&requestLog.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return requestLog, nil
}
