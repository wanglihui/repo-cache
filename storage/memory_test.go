package storage_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/wanglihui/repo-cache/storage"
)

var m = storage.NewMemoryStorage("testmemory")

var (
	id  = "#1"
	val = "test"
)

func TestMemoryStorage_Set(t *testing.T) {
	type args struct {
		ctx context.Context
		key storage.Key
		val storage.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set",
			args: args{
				ctx: context.Background(),
				key: storage.Key(id),
				val: []byte(val),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.Set(tt.args.ctx, tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStorage_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key storage.Key
	}
	tests := []struct {
		name    string
		args    args
		want    storage.Value
		wantErr bool
	}{
		{
			name: "get#1",
			args: args{
				ctx: context.Background(),
				key: storage.Key(id),
			},
			want:    storage.Value(val),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		key storage.Key
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete#1",
			args: args{
				ctx: context.Background(),
				key: storage.Key(id),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.Delete(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
