package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	image_events "github.com/Sokol111/ecommerce-image-service-api/gen/events"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"github.com/Sokol111/ecommerce-product-query-service/internal/infrastructure/client"
	product_events "github.com/Sokol111/ecommerce-product-service-api/gen/events"
	"go.uber.org/zap"
)

type productHandler struct {
	repo       productview.Repository
	attrClient client.AttributeClient
}

func newProductHandler(repo productview.Repository, attrClient client.AttributeClient) *productHandler {
	return &productHandler{
		repo:       repo,
		attrClient: attrClient,
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
	attributes, attrs, err := h.enrichAndMapAttributes(ctx, e.Payload.Attributes)
	if err != nil {
		return fmt.Errorf("failed to enrich attributes: %w", err)
	}

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

	h.log(ctx).Debug("product view created",
		zap.String("productID", e.Payload.ProductID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func (h *productHandler) handleProductUpdated(ctx context.Context, e *product_events.ProductUpdatedEvent) error {
	attributes, attrs, err := h.enrichAndMapAttributes(ctx, e.Payload.Attributes)
	if err != nil {
		return fmt.Errorf("failed to enrich attributes: %w", err)
	}

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

// enrichAndMapAttributes fetches attribute slugs from attribute-service and builds domain attributes
func (h *productHandler) enrichAndMapAttributes(ctx context.Context, eventAttrs *[]product_events.ProductAttribute) ([]productview.ProductAttribute, map[string]any, error) {
	if eventAttrs == nil || len(*eventAttrs) == 0 {
		return nil, nil, nil
	}

	// Collect attribute IDs
	ids := make([]string, len(*eventAttrs))
	for i, attr := range *eventAttrs {
		ids[i] = attr.AttributeID
	}

	// Fetch attribute data from attribute-service
	attrDataMap, err := h.attrClient.GetAttributesByIDs(ctx, ids)
	if err != nil {
		return nil, nil, err
	}

	attributes := make([]productview.ProductAttribute, len(*eventAttrs))
	attrs := make(map[string]any, len(*eventAttrs))

	for i, attr := range *eventAttrs {
		attrData := attrDataMap[attr.AttributeID]
		slug := attrData.Slug

		var optionSlugValues []string
		if attr.OptionSlugValues != nil {
			optionSlugValues = *attr.OptionSlugValues
		}

		attributes[i] = productview.ProductAttribute{
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

	return attributes, attrs, nil
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
