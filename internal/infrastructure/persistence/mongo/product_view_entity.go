package mongo

import (
	"time"
)

// productViewEntity represents the MongoDB document structure for product views
type productViewEntity struct {
	ID          string    `bson:"_id"`
	Version     int       `bson:"version"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Price       float32   `bson:"price"`
	Quantity    int       `bson:"quantity"`
	ImageID     *string   `bson:"imageId,omitempty"`
	Enabled     bool      `bson:"enabled"`
	CreatedAt   time.Time `bson:"createdAt"`
	ModifiedAt  time.Time `bson:"modifiedAt"`
}
