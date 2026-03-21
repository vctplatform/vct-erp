package repository

import "context"

// Reader exposes a single-entity lookup contract.
type Reader[T any, ID comparable] interface {
	GetByID(ctx context.Context, id ID) (T, error)
}

// BulkReader exposes a multi-entity lookup contract.
type BulkReader[T any, ID comparable] interface {
	GetByIDs(ctx context.Context, ids []ID) (map[ID]T, error)
}

// Writer exposes create and update operations for an entity.
type Writer[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
}

// Deleter exposes entity deletion.
type Deleter[ID comparable] interface {
	Delete(ctx context.Context, id ID) error
}

// Repository composes the generic CRUD contracts used across modules.
type Repository[T any, ID comparable] interface {
	Reader[T, ID]
	BulkReader[T, ID]
	Writer[T]
	Deleter[ID]
}

// IsolationLevel declares the transaction isolation requested by a use case.
type IsolationLevel string

const (
	IsolationReadCommitted  IsolationLevel = "read_committed"
	IsolationRepeatableRead IsolationLevel = "repeatable_read"
	IsolationSerializable   IsolationLevel = "serializable"
)

// TxOptions controls the transaction behavior.
type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

// TxManager runs a function inside a transaction boundary.
type TxManager interface {
	WithinTransaction(ctx context.Context, opts TxOptions, fn func(ctx context.Context) error) error
}
