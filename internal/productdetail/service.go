package categorylist

import (
	"context"

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
	return s.store.Upsert(ctx, e.Payload.ProductID, e.Payload.Name, e.Payload.Version, e.Payload.Enabled)
}

func (s *service) ProcessProductUpdatedEvent(ctx context.Context, e *event.Event[payload.ProductUpdated]) error {
	return s.store.Upsert(ctx, e.Payload.ProductID, e.Payload.Name, e.Payload.Version, e.Payload.Enabled)
}

func (s *service) GetById(ctx context.Context, id string) (*model.ProductDTO, error) {
	return s.store.GetAllEnabled(ctx)
}
