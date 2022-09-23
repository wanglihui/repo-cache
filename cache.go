package repocache

import (
	"context"
	"fmt"

	"github.com/wanglihui/repo-cache/storage"
	"golang.org/x/sync/errgroup"
)

type RepoCacheInterface[T EntityModelInterface] interface {
	FindByID(context.Context, ID) (T, error)
	Update(context.Context, T) (T, error)
	Delete(context.Context, ID) error
	Create(context.Context, T) (T, error)
	Find(context.Context, Where, Order, Limit) (Paginate[T], error)
}

type RepoCache[T EntityModelInterface] struct {
	repo    RepoInterface[T]
	storage storage.StorageInterface
}

func (it *RepoCache[T]) FindByID(ctx context.Context, id ID) (T, error) {
	bs, err := it.storage.Get(ctx, storage.Key(id))
	var m T
	if err != nil && bs != nil {
		m.Deserialize(bs)
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

func (it *RepoCache[T]) Update(ctx context.Context, m T) (T, error) {
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

func (it *RepoCache[T]) Delete(ctx context.Context, id ID) error {
	if err := it.repo.Delete(ctx, id); err != nil {
		return err
	}
	return it.storage.Delete(ctx, storage.Key(id))
}

func (it *RepoCache[T]) Create(ctx context.Context, m T) (T, error) {
	m, err := it.repo.Create(ctx, m)
	if err != nil {
		return m, err
	}
	err = it.storage.Set(ctx, storage.Key(m.GetID()), m.Serialize())
	return m, err
}

func (it *RepoCache[T]) Find(ctx context.Context, where Where, order Order, limit Limit) (Paginate[T], error) {
	paginate, err := it.repo.Find(ctx, where, order, limit)
	var p Paginate[T]
	if err != nil {
		return p, err
	}
	items := make([]T, len(paginate.Items))
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

func NewRepoCache[T EntityModelInterface](repo RepoInterface[T], storage storage.StorageInterface) RepoCacheInterface[T] {
	return &RepoCache[T]{
		repo:    repo,
		storage: storage,
	}
}
