package product

import (
	"context"
)

type Service interface {

	// can return NotFoundError
	GetById(ctx context.Context, id string) (*Product, error)

	GetAll(ctx context.Context) ([]*Product, error)
}
