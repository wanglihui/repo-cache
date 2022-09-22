package gormcache

import (
	"context"
	"fmt"

	"github.com/wanglihui/repo-cache/storage"
	"golang.org/x/sync/errgroup"
)

type ReopCacheInterface interface {
	FindByID(context.Context, ID) (EntityModelInterface, error)
	Update(context.Context, EntityModelInterface) (EntityModelInterface, error)
	Delete(context.Context, ID) error
	Create(context.Context, EntityModelInterface) (EntityModelInterface, error)
	Find(context.Context, Where, Order, Limit) (Paginate, error)
}

type RepoCache struct {
	repo    RepoInterface
	storage storage.StorageInterface
}

func (it *RepoCache) FindByID(ctx context.Context, id ID) (EntityModelInterface, error) {
	bs, err := it.storage.Get(ctx, storage.Key(id))
	var m EntityModelInterface
	if err != nil && bs != nil {
		m = m.Deserialize(bs)
		return m, nil
	}
	m, err = it.repo.FindByID(ctx, id)
	if err != nil {
		return m, err
	}
	if m.GetID() == id {
		if err := it.storage.Set(ctx, storage.Key(m.GetID()), m.Serialize()); err != nil {
			fmt.Println(err)
		}
	}
	return m, err
}

func (it *RepoCache) Update(ctx context.Context, m EntityModelInterface) (EntityModelInterface, error) {
	m, err := it.repo.Update(ctx, m)
	if err != nil {
		return m, err
	}
	key := storage.Key(m.GetID())
	if err := it.storage.Delete(ctx, key); err != nil {
		return m, err
	}
	err = it.storage.Set(ctx, key, m.Serialize())
	return m, err
}

func (it *RepoCache) Delete(ctx context.Context, id ID) error {
	if err := it.repo.Delete(ctx, id); err != nil {
		return err
	}
	return it.storage.Delete(ctx, storage.Key(id))
}

func (it *RepoCache) Create(ctx context.Context, m EntityModelInterface) (EntityModelInterface, error) {
	m, err := it.repo.Create(ctx, m)
	if err != nil {
		return m, err
	}
	err = it.storage.Set(ctx, storage.Key(m.GetID()), m.Serialize())
	return m, err
}

func (it *RepoCache) Find(ctx context.Context, where Where, order Order, limit Limit) (Paginate, error) {
	paginate, err := it.repo.Find(ctx, where, order, limit)
	var p Paginate
	if err != nil {
		return p, err
	}
	items := make([]EntityModelInterface, len(paginate.Items))
	var group errgroup.Group
	for idx, id := range paginate.Items {
		idx, id := idx, id
		group.Go(func() error {
			m, err := it.FindByID(ctx, id)
			if err != nil {
				return err
			}
			items[idx] = m
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return p, err
	}
	p.Items = items
	p.Limit = paginate.Limit
	p.Total = paginate.Total
	return p, err
}

func NewGormCache(repo RepoInterface, storage storage.StorageInterface) ReopCacheInterface {
	return &RepoCache{
		repo:    repo,
		storage: storage,
	}
}
