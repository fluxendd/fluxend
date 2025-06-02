package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	flxErrs "fluxend/pkg/errors"
	"fmt"

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

// Common convenience methods for the adapter //

// SelectList executes a query and scans the results into a slice
func (a *Adapter) SelectList(dest interface{}, query string, args ...interface{}) error {
	rows, err := a.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return a.scanRowsIntoSlice(dest, rows)
}

// SelectNamedList executes a named query and scans the results into a slice
func (a *Adapter) SelectNamedList(dest interface{}, query string, arg interface{}) error {
	rows, err := a.db.NamedQuery(query, arg)
	if err != nil {
		return err
	}
	defer rows.Close()

	return a.scanSqlxRowsIntoSlice(dest, rows)
}

// GetWithNotFound executes a Get query with automatic NotFound error handling
func (a *Adapter) GetWithNotFound(dest interface{}, notFoundMsg string, query string, args ...interface{}) error {
	err := a.db.Get(dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return flxErrs.NewNotFoundError(notFoundMsg)
		}

		return err
	}

	return nil
}

// Exists executes an EXISTS query and returns the boolean result
func (a *Adapter) Exists(table, condition string, args ...interface{}) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", table, condition)

	var exists bool
	err := a.db.Get(&exists, query, args...)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// ExecWithRowsAffected executes a query and returns the number of rows affected
func (a *Adapter) ExecWithRowsAffected(query string, args ...interface{}) (int64, error) {
	res, err := a.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// ExecWithErr executes a query and returns only if error occurs
func (a *Adapter) ExecWithErr(query string, args ...interface{}) error {
	_, err := a.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// NamedExecWithRowsAffected executes a named query and returns the number of rows affected
func (a *Adapter) NamedExecWithRowsAffected(query string, arg interface{}) (int64, error) {
	res, err := a.db.NamedExec(query, arg)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// WithTransaction executes a function within a transaction
func (a *Adapter) WithTransaction(fn func(tx shared.Tx) error) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	txAdapter := NewTxAdapter(tx)

	if err := fn(txAdapter); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// === Helper methods for scanning ===

// scanRowsIntoSlice scans regular sql.Rows into a slice using reflection
func (a *Adapter) scanRowsIntoSlice(dest interface{}, rows *sql.Rows) error {
	// This would need to use reflection to dynamically scan into the slice
	// For now, we'll use sqlx's functionality by converting to sqlx.Rows
	sqlxRows := &sqlx.Rows{Rows: rows}

	return a.scanSqlxRowsIntoSlice(dest, sqlxRows)
}

// scanSqlxRowsIntoSlice scans sqlx.Rows into a slice using StructScan
func (a *Adapter) scanSqlxRowsIntoSlice(dest interface{}, rows *sqlx.Rows) error {
	// Use sqlx's built-in StructScan functionality
	// This assumes dest is a pointer to a slice of structs
	err := sqlx.StructScan(rows, dest)
	if err != nil {
		return pkg.FormatError(err, "scan", pkg.GetMethodName())
	}

	if err := rows.Err(); err != nil {
		return pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return nil
}
