package productview

import (
	"context"

	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

// AttributeFilter represents a filter for a single attribute
type AttributeFilter struct {
	Slug   string
	Values []string // For single/multiple type attributes
	Min    *float64 // For range type attributes
	Max    *float64 // For range type attributes
}

type ListQuery struct {
	Page             int
	Size             int
	CategoryID       *string
	Sort             string
	Order            string
	MinPrice         *float64
	MaxPrice         *float64
	AttributeFilters []AttributeFilter
}

type Repository interface {
	// Upsert inserts or updates a product view (for event processing)
	Upsert(ctx context.Context, product *ProductView) error

	// UpdateImageURL updates the image URL for a product (called when ImagePromoted event is received)
	// Returns nil if product doesn't exist yet (image event arrived before product event)
	UpdateImageURL(ctx context.Context, productID, imageID, imageURL string) error

	// FindByID retrieves a product view by ID
	FindByID(ctx context.Context, id string) (*ProductView, error)

	// FindRandom retrieves random enabled products
	FindRandom(ctx context.Context, count int) ([]*ProductView, error)

	// FindList retrieves a paginated list of products
	FindList(ctx context.Context, query ListQuery) (*commonsmongo.PageResult[ProductView], error)
}
