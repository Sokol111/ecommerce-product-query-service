package productview

import "time"

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
	CategoryID  *string
	Enabled     bool
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

// Reconstruct creates a ProductView from persistence data
func Reconstruct(id string, version int, name string, description *string, price float32, quantity int, imageID *string, categoryID *string, enabled bool, createdAt, modifiedAt time.Time) *ProductView {
	return &ProductView{
		ID:          id,
		Version:     version,
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		ImageID:     imageID,
		CategoryID:  categoryID,
		Enabled:     enabled,
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
	}
}

// NewProductView creates a new product view from event data
func NewProductView(id string, version int, name string, description *string, price float32, quantity int, imageID *string, categoryID *string, enabled bool, createdAt, modifiedAt time.Time) *ProductView {
	return &ProductView{
		ID:          id,
		Version:     version,
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		ImageID:     imageID,
		CategoryID:  categoryID,
		Enabled:     enabled,
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
	}
}
