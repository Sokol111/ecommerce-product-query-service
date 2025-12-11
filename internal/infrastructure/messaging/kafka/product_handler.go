package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"github.com/Sokol111/ecommerce-product-service-api/gen/events"
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
	// Type assert to Event interface first to get exhaustiveness checking
	e, ok := event.(events.Event)
	if !ok {
		return fmt.Errorf("event does not implement Event interface: %T: %w", event, consumer.ErrSkipMessage)
	}

	// Now switch on concrete types - exhaustive linter will warn if any Event type is missing
	switch evt := e.(type) {
	case *events.ProductCreatedEvent:
		return h.handleProductCreated(ctx, evt)
	case *events.ProductUpdatedEvent:
		return h.handleProductUpdated(ctx, evt)
	default:
		// If exhaustive linter is enabled and all Event types are handled above,
		// this case should theoretically never be reached
		return fmt.Errorf("unhandled event type: %T: %w", event, consumer.ErrSkipMessage)
	}
}

func (h *productHandler) handleProductCreated(ctx context.Context, e *events.ProductCreatedEvent) error {
	view := productview.NewProductView(
		e.Payload.ProductID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Description,
		e.Payload.Price,
		e.Payload.Quantity,
		e.Payload.ImageID,
		e.Payload.CategoryID,
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
	view := productview.NewProductView(
		e.Payload.ProductID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Description,
		e.Payload.Price,
		e.Payload.Quantity,
		e.Payload.ImageID,
		e.Payload.CategoryID,
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
	return logger.Get(ctx).With(zap.String("component", "product-handler"))
}
