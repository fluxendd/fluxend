package sqlx

import (
	"context"
	"database/sql"
	"fluxend/internal/domain/shared"

	"github.com/jmoiron/sqlx"
)

// TxAdapter wraps sqlx.Tx and implements the shared.Tx interface
type TxAdapter struct {
	tx *sqlx.Tx
}

// NewTxAdapter creates a new transaction adapter
func NewTxAdapter(tx *sqlx.Tx) shared.Tx {
	return &TxAdapter{
		tx: tx,
	}
}

func (t *TxAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

func (t *TxAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *TxAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

func (t *TxAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *TxAdapter) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return t.tx.QueryRowx(query, args...)
}

func (t *TxAdapter) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return t.tx.NamedQuery(query, arg)
}

func (t *TxAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *TxAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *TxAdapter) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return t.tx.NamedExec(query, arg)
}

func (t *TxAdapter) Get(dest interface{}, query string, args ...interface{}) error {
	return t.tx.Get(dest, query, args...)
}

func (t *TxAdapter) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.GetContext(ctx, dest, query, args...)
}

func (t *TxAdapter) Select(dest interface{}, query string, args ...interface{}) error {
	return t.tx.Select(dest, query, args...)
}

func (t *TxAdapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

func (t *TxAdapter) Commit() error {
	return t.tx.Commit()
}

func (t *TxAdapter) Rollback() error {
	return t.tx.Rollback()
}
