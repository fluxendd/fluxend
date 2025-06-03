package shared

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DatabaseService interface {
	Create(name string, userUUID uuid.NullUUID) error
	DropIfExists(name string) error
	Recreate(name string) error
	List() ([]string, error)
	Exists(name string) (bool, error)
	Connect(name string) (*sqlx.DB, error)
}

type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Beginx() (Tx, error)

	// Custom convenience methods
	SelectNamedList(dest interface{}, query string, arg interface{}) error
	GetWithNotFound(dest interface{}, notFoundMsg string, query string, args ...interface{}) error
	Exists(table, condition string, args ...interface{}) (bool, error)
	ExecWithErr(query string, args ...interface{}) error
	ExecWithRowsAffected(query string, args ...interface{}) (int64, error)
	NamedExecWithRowsAffected(query string, arg interface{}) (int64, error)
	WithTransaction(fn func(tx Tx) error) error
}

type Tx interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	//NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)

	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Commit() error
	Rollback() error
}
