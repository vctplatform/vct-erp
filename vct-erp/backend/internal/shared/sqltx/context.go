package sqltx

import (
	"context"
	"database/sql"
)

type txContextKey struct{}

// WithTx stores a SQL transaction in context for nested repository calls.
func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// FromContext returns the SQL transaction stored in context, if any.
func FromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(*sql.Tx)
	return tx, ok && tx != nil
}
