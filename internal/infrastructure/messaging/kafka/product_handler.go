package kafka

import (
	"context"
	"fmt"

	catalog_events "github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	image_events "github.com/Sokol111/ecommerce-image-service-api/gen/events"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
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
	case *catalog_events.ProductUpdatedEvent:
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

func (h *productHandler) handleProductUpdated(ctx context.Context, e *catalog_events.ProductUpdatedEvent) error {
	attributes, attrs := mapAttributes(e.Payload.Attributes)

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
		attributes,
		attrs,
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

// mapAttributes converts event attributes to domain attributes.
// Only immutable fields (IDs, slugs) and product-specific values are mapped.
// Mutable display data will be joined from attributes collection at read time.
func mapAttributes(eventAttrs *[]catalog_events.AttributeValue) ([]productview.AttributeValue, map[string]any) {
	if eventAttrs == nil || len(*eventAttrs) == 0 {
		return nil, nil
	}

	attributes := make([]productview.AttributeValue, len(*eventAttrs))
	attrs := make(map[string]any, len(*eventAttrs))

	for i, attr := range *eventAttrs {
		slug := attr.AttributeSlug

		var optionSlugValues []string
		if attr.OptionSlugValues != nil {
			optionSlugValues = *attr.OptionSlugValues
		}

		attributes[i] = productview.AttributeValue{
			AttributeID:      attr.AttributeID,
			Slug:             slug,
			OptionSlugValue:  attr.OptionSlugValue,
			OptionSlugValues: optionSlugValues,
			NumericValue:     attr.NumericValue,
			TextValue:        attr.TextValue,
			BooleanValue:     attr.BooleanValue,
		}

		// Build attrs map for filtering
		if slug != "" {
			if attr.NumericValue != nil {
				attrs[slug] = *attr.NumericValue
			} else if attr.OptionSlugValues != nil && len(*attr.OptionSlugValues) > 0 {
				attrs[slug] = *attr.OptionSlugValues
			} else if attr.OptionSlugValue != nil {
				attrs[slug] = *attr.OptionSlugValue
			} else if attr.TextValue != nil {
				attrs[slug] = *attr.TextValue
			} else if attr.BooleanValue != nil {
				attrs[slug] = *attr.BooleanValue
			}
		}
	}

	return attributes, attrs
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
