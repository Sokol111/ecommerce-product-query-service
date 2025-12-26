package mongo

import (
	"time"
)

// productAttributeEntity represents the MongoDB subdocument structure for product attributes
type productAttributeEntity struct {
	AttributeID  string   `bson:"attributeId"`
	Value        *string  `bson:"value,omitempty"`
	Values       []string `bson:"values,omitempty"`
	NumericValue *float32 `bson:"numericValue,omitempty"`
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
}
