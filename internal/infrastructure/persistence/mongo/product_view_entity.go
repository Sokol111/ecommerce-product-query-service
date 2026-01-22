package mongo

import (
	"time"
)

// productAttributeEntity represents the MongoDB subdocument structure for product attributes.
// Only immutable fields (IDs, slugs) and product-specific values are stored.
// Mutable display data should be joined from attributes collection.
type productAttributeEntity struct {
	AttributeID      string   `bson:"attributeId"`
	Slug             string   `bson:"slug"`
	OptionSlugValue  *string  `bson:"optionSlugValue,omitempty"`
	OptionSlugValues []string `bson:"optionSlugValues,omitempty"`
	NumericValue     *float32 `bson:"numericValue,omitempty"`
	TextValue        *string  `bson:"textValue,omitempty"`
	BooleanValue     *bool    `bson:"booleanValue,omitempty"`
}

// productViewEntity represents the MongoDB document structure for product views
type productViewEntity struct {
	ID          string                   `bson:"_id"`
	Version     int                      `bson:"version"`
	Name        string                   `bson:"name"`
	Description *string                  `bson:"description,omitempty"`
	Price       float32                  `bson:"price"`
	Quantity    int                      `bson:"quantity"`
	ImageID     *string                  `bson:"imageId,omitempty"`
	ImageURL    *string                  `bson:"imageUrl,omitempty"`
	CategoryID  *string                  `bson:"categoryId,omitempty"`
	Enabled     bool                     `bson:"enabled"`
	CreatedAt   time.Time                `bson:"createdAt"`
	ModifiedAt  time.Time                `bson:"modifiedAt"`
	Attributes  []productAttributeEntity `bson:"attributes,omitempty"`
	Attrs       map[string]any           `bson:"attrs,omitempty"` // Denormalized for filtering
}
