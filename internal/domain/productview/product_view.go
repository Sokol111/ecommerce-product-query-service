package productview

import "time"

// AttributeType represents the type of attribute
type AttributeType string

const (
	AttributeTypeSingle   AttributeType = "single"
	AttributeTypeMultiple AttributeType = "multiple"
	AttributeTypeRange    AttributeType = "range"
	AttributeTypeBoolean  AttributeType = "boolean"
	AttributeTypeText     AttributeType = "text"
)

// AttributeRole defines how an attribute is used in a category
type AttributeRole string

const (
	AttributeRoleVariant       AttributeRole = "variant"
	AttributeRoleSpecification AttributeRole = "specification"
)

// ProductAttribute represents product attribute with value and display information
type ProductAttribute struct {
	AttributeID      string
	Slug             string
	Name             string        // Display name of the attribute
	Type             AttributeType // Attribute type for rendering
	Unit             *string       // Unit for range type attributes
	Role             AttributeRole // How attribute is used (variant/specification)
	SortOrder        int           // Display order of the attribute
	OptionSlugValue  *string
	OptionSlugValues []string
	OptionName       *string  // Display name of selected option (for single type)
	OptionNames      []string // Display names of selected options (for multiple type)
	OptionColorCode  *string  // Hex color code for color attributes
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
