package productview

import "time"

// ProductAttribute represents product attribute with value
type ProductAttribute struct {
	AttributeID      string
	Slug             string
	OptionSlugValue  *string
	OptionSlugValues []string
	NumericValue     *float32
	TextValue        *string
	BooleanValue     *bool
}

// ProductView - read model for product queries (CQRS query side)
// Unlike the domain Product in the command service, this is a denormalized view optimized for reads
type ProductView struct {
	ID          string
	Version     int
	Name        string
	Description *string
	Price       float32
	Quantity    int
	ImageID     *string
	ImageURL    *string // Denormalized image URL for efficient reads
	CategoryID  *string
	Enabled     bool
	CreatedAt   time.Time
	ModifiedAt  time.Time
	Attributes  []ProductAttribute
	Attrs       map[string]any // Denormalized attributes map for filtering (slug -> value)
}

// Reconstruct creates a ProductView from persistence data
func Reconstruct(id string, version int, name string, description *string, price float32, quantity int, imageID *string, imageURL *string, categoryID *string, enabled bool, createdAt, modifiedAt time.Time, attributes []ProductAttribute, attrs map[string]any) *ProductView {
	return &ProductView{
		ID:          id,
		Version:     version,
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		ImageID:     imageID,
		ImageURL:    imageURL,
		CategoryID:  categoryID,
		Enabled:     enabled,
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
		Attributes:  attributes,
		Attrs:       attrs,
	}
}

// NewProductView creates a new product view from event data
func NewProductView(id string, version int, name string, description *string, price float32, quantity int, imageID *string, categoryID *string, enabled bool, createdAt, modifiedAt time.Time, attributes []ProductAttribute, attrs map[string]any) *ProductView {
	return &ProductView{
		ID:          id,
		Version:     version,
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		ImageID:     imageID,
		ImageURL:    nil, // Will be populated by ImagePromoted event
		CategoryID:  categoryID,
		Enabled:     enabled,
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
		Attributes:  attributes,
		Attrs:       attrs,
	}
}
