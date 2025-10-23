package productdetail

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-product-query-service/internal/model"
)

type service struct {
	store Store
}

func newService(store Store) model.ProductDetailService {
	return &service{store: store}
}

func (s *service) ProcessProductCreatedEvent(ctx context.Context, e *event.Event[payload.ProductCreated]) error {
	return s.store.Upsert(ctx, &model.Product{
		ID:          e.Payload.ProductID,
		Name:        e.Payload.Name,
		Description: e.Payload.Description,
		Enabled:     e.Payload.Enabled,
		Quantity:    e.Payload.Quantity,
		Price:       e.Payload.Price,
		ImageId:     e.Payload.ImageId,
		CreatedAt:   e.Payload.CreatedAt,
		ModifiedAt:  e.Payload.ModifiedAt,
		Version:     e.Payload.Version,
	})
}

func (s *service) ProcessProductUpdatedEvent(ctx context.Context, e *event.Event[payload.ProductUpdated]) error {
	return s.store.Upsert(ctx, &model.Product{
		ID:          e.Payload.ProductID,
		Name:        e.Payload.Name,
		Description: e.Payload.Description,
		Enabled:     e.Payload.Enabled,
		Quantity:    e.Payload.Quantity,
		Price:       e.Payload.Price,
		ImageId:     e.Payload.ImageId,
		CreatedAt:   e.Payload.CreatedAt,
		ModifiedAt:  e.Payload.ModifiedAt,
		Version:     e.Payload.Version,
	})
}

func (s *service) GetRandomProducts(ctx context.Context, amount int) ([]*model.ProductDTO, error) {
	products, err := s.store.GetRandomProducts(ctx, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to get random products, amount [%v]: %w", amount, err)
	}
	return products, nil
}

func (s *service) GetById(ctx context.Context, id string) (*model.ProductDTO, error) {
	p, err := s.store.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errEntityNotFound) {
			return nil, fmt.Errorf("product not found, id [%v]: %w", id, model.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get product by id [%v]: %w", id, err)
	}
	return p, nil
}
