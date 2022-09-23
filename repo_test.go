package repocache_test

import (
	"context"
	"fmt"
	"time"

	repocache "github.com/wanglihui/repo-cache"
)

func NewRepoImpl() *RepoImpl {
	return &RepoImpl{
		db: make(map[string]EntityTest, 0),
	}
}

type RepoImpl struct {
	db map[string]EntityTest
}

func (r *RepoImpl) FindByID(ctx context.Context, id repocache.ID) (EntityTest, error) {
	m := r.db[string(id)]
	return m, nil
}

func (r *RepoImpl) Update(ctx context.Context, m EntityTest) (EntityTest, error) {
	id := m.GetID()
	r.db[string(id)] = m
	return m, nil
}

func (r *RepoImpl) Delete(ctx context.Context, id repocache.ID) error {
	delete(r.db, string(id))
	return nil
}

func (r *RepoImpl) Create(ctx context.Context, m EntityTest) (EntityTest, error) {
	if m.ID == "" {
		id := fmt.Sprintf("%d", time.Now().Unix())
		m.ID = repocache.ID(id)
	}
	r.db[string(m.ID)] = m
	return m, nil
}

func (r *RepoImpl) Find(ctx context.Context, where repocache.Where, order repocache.Order, limit repocache.Limit) (repocache.PaginateID, error) {
	var (
		total = len(r.db)
		p     = repocache.PaginateID{
			Limit: limit,
			Total: int64(total),
		}
	)
	var items = make([]repocache.ID, 0)
	for _, v := range r.db {
		items = append(items, v.GetID())
	}
	p.Items = items
	return p, nil
}
