package gormcache

import "context"

type EntityModelInterface interface {
	GetID() ID
	Serialize() []byte //对象序列化
	Deserialize([]byte) EntityModelInterface
}

type Where map[string]string
type Limit int64
type Order string

type PaginateID struct {
	Items []ID  `json:"items"`
	Total int64 `json:"total"`
	Limit Limit `json:"limit"`
}

type ID string

type RepoInterface interface {
	FindByID(context.Context, ID) (EntityModelInterface, error)
	Update(context.Context, EntityModelInterface) (EntityModelInterface, error)
	Delete(context.Context, ID) error
	Create(context.Context, EntityModelInterface) (EntityModelInterface, error)
	Find(context.Context, Where, Order, Limit) (PaginateID, error)
}

type Paginate struct {
	Items []EntityModelInterface `json:"items"`
	Total int64                  `json:"total"`
	Limit Limit                  `json:"limit"`
}
