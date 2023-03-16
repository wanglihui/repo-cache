package repocache

import "context"

type EntityModelInterface interface {
	GetID() ID
	Serialize() []byte //对象序列化
	Deserialize([]byte) (interface{}, error)
}

type Where struct {
	SQL    string
	Args   []interface{}
	Limit  int64
	Offset int64
	Order  string
}
type Limit int64
type Order string

type PaginateID struct {
	Items []ID  `json:"items"`
	Total int64 `json:"total"`
	Limit Limit `json:"limit"`
}

type ID string

type RepoInterface[T EntityModelInterface] interface {
	FindByID(context.Context, ID) (T, error)
	Update(context.Context, T) (T, error)
	Delete(context.Context, ID) error
	Create(context.Context, T) (T, error)
	Find(context.Context, Where) (PaginateID, error)
	FindOne(context.Context, Where) (ID, error)
}

type Paginate[T EntityModelInterface] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Limit Limit `json:"limit"`
}
