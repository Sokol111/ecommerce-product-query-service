package productview

import "context"

type Repository interface {
	// Upsert inserts or updates a product view (for event processing)
	Upsert(ctx context.Context, product *ProductView) error

	// FindByID retrieves a product view by ID
	FindByID(ctx context.Context, id string) (*ProductView, error)

	// FindRandom retrieves random enabled products
	FindRandom(ctx context.Context, amount int) ([]*ProductView, error)
}
