package repocache_test

import (
	"context"
	"testing"

	repocache "github.com/wanglihui/repo-cache"
	"github.com/wanglihui/repo-cache/storage"
)

var (
	repo = NewRepoImpl()
	s    = storage.NewMemoryStorage("unit_test_repo")
	it   = repocache.NewRepoCache[EntityTest](repo, s)
)

func TestRepoCache_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		m   EntityTest
	}
	tests := []struct {
		name    string
		args    args
		want    EntityTest
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				m: EntityTest{
					ID:   "1",
					Name: "test",
					Age:  10,
				},
			},
			want: EntityTest{
				Name: "test",
				Age:  10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := it.Create(tt.args.ctx, tt.args.m)
			defer it.Delete(tt.args.ctx, got.GetID())
			if (err != nil) != tt.wantErr {
				t.Errorf("RepoCache.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Age != tt.want.Age {
				t.Errorf("RepoCache.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoCache_FindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  repocache.ID
	}
	tests := []struct {
		name    string
		args    args
		want    EntityTest
		wantErr bool
	}{
		{
			name: "findById",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			want: EntityTest{
				ID:  "1",
				Age: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it.Create(tt.args.ctx, tt.want)
			got, err := it.FindByID(tt.args.ctx, tt.args.id)
			got, err = it.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Age != tt.want.Age {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoCache_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		m   EntityTest
	}
	tests := []struct {
		name    string
		args    args
		want    EntityTest
		wantErr bool
	}{
		{
			name: "TestRepoCache_Update",
			args: args{
				ctx: context.Background(),
				m: EntityTest{
					ID:  "1",
					Age: 20,
				},
			},
			want: EntityTest{
				ID:  "1",
				Age: 20,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := it.Update(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Age != tt.want.Age {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoCache_Find(t *testing.T) {
	type args struct {
		ctx   context.Context
		where repocache.Where
		order repocache.Order
		limit repocache.Limit
	}
	tests := []struct {
		name    string
		args    args
		want    repocache.Paginate[EntityTest]
		wantErr bool
	}{
		{
			name: "find",
			args: args{
				ctx: context.Background(),
				where: repocache.Where{
					SQL:    "id = ? or id = ?",
					Args:   []interface{}{1, 2},
					Limit:  10,
					Order:  "",
					Offset: 0,
				},
			},
			want: repocache.Paginate[EntityTest]{
				Limit: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := it.Find(tt.args.ctx, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Limit != tt.want.Limit {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if len(got.Items) < 0 {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoCache_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  repocache.ID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestRepoCache_Delete",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := it.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
