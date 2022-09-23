package storage

import (
	"context"
	"fmt"
	"sync"
)

var storages = make(map[string][]byte, 100)

func NewMemoryStorage(prefix string) StorageInterface {
	return &MemoryStorage{
		prefix: prefix,
	}
}

type MemoryStorage struct {
	prefix string
	l      sync.Mutex
}

func (m *MemoryStorage) Get(ctx context.Context, key Key) (Value, error) {
	fullKey := m.getKey(key)
	val := storages[fullKey]
	return val, nil
}

func (m *MemoryStorage) Set(ctx context.Context, key Key, val Value) error {
	fullKey := m.getKey(key)
	m.l.Lock()
	storages[fullKey] = val
	m.l.Unlock()
	return nil
}

func (m *MemoryStorage) Delete(ctx context.Context, key Key) error {
	fullKey := m.getKey(key)
	storages[fullKey] = nil
	return nil
}

func (m *MemoryStorage) getKey(key Key) string {
	return fmt.Sprintf("%s:%s", m.prefix, string(key))
}
