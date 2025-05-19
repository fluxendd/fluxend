package client

import (
	"github.com/jmoiron/sqlx"
)

type ConnectionService interface {
	ConnectByDatabaseName(name string) (*sqlx.DB, error)
	GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error)
	GetTableRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error)
	GetColumnRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error)
	GetIndexRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error)
	GetRowRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error)
}
