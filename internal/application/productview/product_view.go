package productview

import "time"

// AttributeValue represents product attribute reference with value.
// Only immutable fields (IDs, slugs) and product-specific values are stored.
// Mutable display data (name, unit, color_code) should be joined from attributes collection.
type AttributeValue struct {
	AttributeID      string
	Slug             string // Immutable attribute slug for filtering
	OptionSlugValue  *string
	OptionSlugValues []string
	NumericValue     *float64
	TextValue        *string
	BooleanValue     *bool
}

// ProductView - read model for product queries (CQRS query side)
// Unlike the domain Product in the command service, this is a denormalized view optimized for reads
type ProductView struct {
	ID            string
	Version       int
	Name          string
	Description   *string
	Price         float64
	Quantity      int
	ImageID       *string
	SmallImageURL *string // Small image URL (400px) for thumbnails and listings
	LargeImageURL *string // Large image URL (800px) for product detail pages
	CategoryID    *string
	Enabled       bool
	CreatedAt     time.Time
	ModifiedAt    time.Time
	Attributes    []AttributeValue
	Attrs         map[string]any // Denormalized attributes map for filtering (slug -> value)
}

// Reconstruct creates a ProductView from persistence data
func Reconstruct(id string, version int, name string, description *string, price float64, quantity int, imageID *string, smallImageURL *string, largeImageURL *string, categoryID *string, enabled bool, createdAt, modifiedAt time.Time, attributes []AttributeValue, attrs map[string]any) *ProductView {
	return &ProductView{
		ID:            id,
		Version:       version,
		Name:          name,
		Description:   description,
		Price:         price,
		Quantity:      quantity,
		ImageID:       imageID,
		SmallImageURL: smallImageURL,
		LargeImageURL: largeImageURL,
		CategoryID:    categoryID,
		Enabled:       enabled,
		CreatedAt:     createdAt,
		ModifiedAt:    modifiedAt,
		Attributes:    attributes,
		Attrs:         attrs,
	}
}

// NewProductView creates a new product view from event data
func NewProductView(id string, version int, name string, description *string, price float64, quantity int, imageID *string, categoryID *string, enabled bool, createdAt, modifiedAt time.Time, attributes []AttributeValue, attrs map[string]any) *ProductView {
	return &ProductView{
		ID:            id,
		Version:       version,
		Name:          name,
		Description:   description,
		Price:         price,
		Quantity:      quantity,
		ImageID:       imageID,
		SmallImageURL: nil, // Will be populated by ImagePromoted event
		LargeImageURL: nil, // Will be populated by ImagePromoted event
		CategoryID:    categoryID,
		Enabled:       enabled,
		CreatedAt:     createdAt,
		ModifiedAt:    modifiedAt,
		Attributes:    attributes,
		Attrs:         attrs,
	}
}
