// Package model is where we keep our database tables representations as
// structures.
package model

type Table interface {
	Migrate() error
	Create(data Data) (*Data, error)
	All(where string, limit int, offset int) ([]Data, error)
	GetBy(where string) (*Data, error)
	GetById(id int64) (*Data, error)
	Update(id int64, updated Data) (*Data, error)
	Delete(id int64) error
}

type Data struct {
	ID int64 `db:id`
}
