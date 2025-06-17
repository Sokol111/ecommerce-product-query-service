package model

import (
	"context"
	"errors"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
)

var ErrNotFound = errors.New("not found")

type ProductDTO struct {
	ID       string
	Name     string
	Price    float32
	Quantity int
}

type ProductDetailService interface {
	ProcessProductCreatedEvent(ctx context.Context, e *event.Event[payload.ProductCreated]) error

	ProcessProductUpdatedEvent(ctx context.Context, e *event.Event[payload.ProductUpdated]) error

	// can return ErrNotFound
	GetById(ctx context.Context, id string) (*ProductDTO, error)
}
