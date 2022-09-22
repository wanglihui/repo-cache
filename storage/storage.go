package storage

import "context"

type Key string
type Value []byte

type StorageInterface interface {
	Get(context.Context, Key) (Value, error)
	Set(context.Context, Key, Value) error
	Delete(context.Context, Key) error
}
