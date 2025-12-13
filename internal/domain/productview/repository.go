package productview

import (
	"context"

	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

type ListQuery struct {
	Page       int
	Size       int
	CategoryID *string
	Sort       string
	Order      string
}

type Repository interface {
	// Upsert inserts or updates a product view (for event processing)
	Upsert(ctx context.Context, product *ProductView) error

	// FindByID retrieves a product view by ID
	FindByID(ctx context.Context, id string) (*ProductView, error)

	// FindRandom retrieves random enabled products
	FindRandom(ctx context.Context, count int) ([]*ProductView, error)

	// FindList retrieves a paginated list of products
	FindList(ctx context.Context, query ListQuery) (*commonsmongo.PageResult[ProductView], error)
}
