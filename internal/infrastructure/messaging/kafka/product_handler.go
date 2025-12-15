package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	image_events "github.com/Sokol111/ecommerce-image-service-api/gen/events"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	product_events "github.com/Sokol111/ecommerce-product-service-api/gen/events"
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
	switch evt := event.(type) {
	// Product events
	case *product_events.ProductCreatedEvent:
		return h.handleProductCreated(ctx, evt)
	case *product_events.ProductUpdatedEvent:
		return h.handleProductUpdated(ctx, evt)

	// Image events (published to same topic with product_id as partition key)
	case *image_events.ProductImagePromotedEvent:
		return h.handleProductImagePromoted(ctx, evt)

	default:
		logger.Get(ctx).Warn("unknown event type, skipping",
			zap.String("type", fmt.Sprintf("%T", event)))
		return fmt.Errorf("unhandled event type: %T: %w", event, consumer.ErrSkipMessage)
	}
}

func (h *productHandler) handleProductCreated(ctx context.Context, e *product_events.ProductCreatedEvent) error {
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

func (h *productHandler) handleProductUpdated(ctx context.Context, e *product_events.ProductUpdatedEvent) error {
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

func (h *productHandler) handleProductImagePromoted(ctx context.Context, e *image_events.ProductImagePromotedEvent) error {
	if err := h.repo.UpdateImageURL(ctx, e.Payload.ProductID, e.Payload.ImageID, e.Payload.ImageURL); err != nil {
		return fmt.Errorf("failed to update product image URL: %w", err)
	}

	h.log(ctx).Debug("product image URL updated",
		zap.String("productID", e.Payload.ProductID),
		zap.String("imageID", e.Payload.ImageID),
		zap.String("eventID", e.Metadata.EventID))

	return nil
}

func (h *productHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "product-handler"))
}
