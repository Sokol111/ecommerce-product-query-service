package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/messaging"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type productUpdatedHandler struct {
	repo productview.Repository
}

func newProductUpdatedHandler(repo productview.Repository) consumer.Handler[messaging.ProductUpdated] {
	return &productUpdatedHandler{
		repo: repo,
	}
}

func (h *productUpdatedHandler) Process(ctx context.Context, e *messaging.Event[messaging.ProductUpdated]) error {
	view := productview.NewProductView(
		e.Payload.ProductID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Description,
		e.Payload.Price,
		e.Payload.Quantity,
		e.Payload.ImageId,
		e.Payload.Enabled,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	return h.repo.Upsert(ctx, view)
}

func (h *productUpdatedHandler) Validate(payload *messaging.ProductUpdated) error {
	return nil
}
