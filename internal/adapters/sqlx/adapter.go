package sqlx

import (
	"context"
	"database/sql"
	"fluxend/internal/domain/shared"

	"github.com/jmoiron/sqlx"
)

// Adapter wraps sqlx.DB and implements the shared.DB interface
// At the moment, it only ports to underlying sqlx.DB methods.
// The purpose of this adapter is to provide a consistent interface
// for the database operations in the application, allowing for easier testing and mocking.
// This adapter can be extended in the future to add more functionality or custom methods.
type Adapter struct {
	db *sqlx.DB
}

// NewAdapter creates a new SQLXAdapter
func NewAdapter(db *sqlx.DB) shared.DB {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.Query(query, args...)
}

func (a *Adapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.QueryContext(ctx, query, args...)
}

func (a *Adapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.db.QueryRow(query, args...)
}

func (a *Adapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return a.db.QueryRowContext(ctx, query, args...)
}

func (a *Adapter) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return a.db.NamedQuery(query, arg)
}

func (a *Adapter) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return a.db.NamedQueryContext(ctx, query, arg)
}

func (a *Adapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.db.Exec(query, args...)
}

func (a *Adapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return a.db.ExecContext(ctx, query, args...)
}

func (a *Adapter) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return a.db.NamedExec(query, arg)
}

func (a *Adapter) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return a.db.NamedExecContext(ctx, query, arg)
}

func (a *Adapter) Get(dest interface{}, query string, args ...interface{}) error {
	return a.db.Get(dest, query, args...)
}

func (a *Adapter) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return a.db.GetContext(ctx, dest, query, args...)
}

func (a *Adapter) Select(dest interface{}, query string, args ...interface{}) error {
	return a.db.Select(dest, query, args...)
}

func (a *Adapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return a.db.SelectContext(ctx, dest, query, args...)
}

func (a *Adapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (shared.Tx, error) {
	tx, err := a.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewTxAdapter(tx), nil
}

func (a *Adapter) Beginx() (shared.Tx, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return nil, err
	}
	return NewTxAdapter(tx), nil
}
