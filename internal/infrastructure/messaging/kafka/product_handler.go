package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"github.com/Sokol111/ecommerce-product-service-api/events"
	"go.uber.org/zap"
)

type productHandler struct {
	repo productview.Repository
}

func newProductHandler(repo productview.Repository) *productHandler {
	return &productHandler{
		repo: repo,
	}
}

func (h *productHandler) Process(ctx context.Context, event any) error {
	switch e := event.(type) {
	case *events.ProductCreatedEvent:
		return h.handleProductCreated(ctx, e)
	case *events.ProductUpdatedEvent:
		return h.handleProductUpdated(ctx, e)
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
}

func (h *productHandler) handleProductCreated(ctx context.Context, e *events.ProductCreatedEvent) error {
	// TODO: Description field is missing in the event schema
	// Consider adding it to product_created.avsc and product_updated.avsc
	view := productview.NewProductView(
		e.Payload.ProductID,
		e.Payload.Version,
		e.Payload.Name,
		"", // description is not included in the current event schema
		e.Payload.Price,
		e.Payload.Quantity,
		e.Payload.ImageID,
		e.Payload.Enabled,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	if err := h.repo.Upsert(ctx, view); err != nil {
		return fmt.Errorf("failed to upsert product view: %w", err)
	}

	h.log(ctx).Debug("product view created",
		zap.String("productID", e.Payload.ProductID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func (h *productHandler) handleProductUpdated(ctx context.Context, e *events.ProductUpdatedEvent) error {
	// TODO: Description field is missing in the event schema
	// Consider adding it to product_created.avsc and product_updated.avsc
	view := productview.NewProductView(
		e.Payload.ProductID,
		e.Payload.Version,
		e.Payload.Name,
		"", // description is not included in the current event schema
		e.Payload.Price,
		e.Payload.Quantity,
		e.Payload.ImageID,
		e.Payload.Enabled,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	if err := h.repo.Upsert(ctx, view); err != nil {
		return fmt.Errorf("failed to upsert product view: %w", err)
	}

	h.log(ctx).Debug("product view updated",
		zap.String("productID", e.Payload.ProductID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func (h *productHandler) log(ctx context.Context) *zap.Logger {
	return logger.FromContext(ctx).With(zap.String("component", "product-handler"))
}
