package client

import (
	"fluxton/internal/domain/database/column"
	"fluxton/internal/domain/database/index"
	"fluxton/internal/domain/database/row"
	"fluxton/internal/domain/database/stat"
	"fluxton/internal/domain/database/table"
	"github.com/jmoiron/sqlx"
)

type ConnectionService interface {
	ConnectByDatabaseName(name string) (*sqlx.DB, error)
	GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (stat.Repository, *sqlx.DB, error)
	GetTableRepo(databaseName string, connection *sqlx.DB) (table.Repository, *sqlx.DB, error)
	GetColumnRepo(databaseName string, connection *sqlx.DB) (column.Repository, *sqlx.DB, error)
	GetIndexRepo(databaseName string, connection *sqlx.DB) (index.Repository, *sqlx.DB, error)
	GetRowRepo(databaseName string, connection *sqlx.DB) (row.Repository, *sqlx.DB, error)
}
