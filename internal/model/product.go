package model

import (
	"context"
	"errors"
	"time"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
)

var ErrNotFound = errors.New("not found")

type Product struct {
	ID          string
	Version     int
	Name        string
	Description string
	Price       float32
	Quantity    int
	ImageId     *string
	Enabled     bool
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

type ProductDTO struct {
	ID       string
	Name     string
	Price    float32
	Quantity int
	ImageId  *string
}

type ProductDetailService interface {
	ProcessProductCreatedEvent(ctx context.Context, e *event.Event[payload.ProductCreated]) error

	ProcessProductUpdatedEvent(ctx context.Context, e *event.Event[payload.ProductUpdated]) error

	// can return ErrNotFound
	GetById(ctx context.Context, id string) (*ProductDTO, error)

	GetRandomProducts(ctx context.Context, amount int) ([]*ProductDTO, error)
}
