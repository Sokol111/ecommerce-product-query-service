package kafka

import (
	"context"

	catalog_events "github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	image_events "github.com/Sokol111/ecommerce-image-service-api/gen/events"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/productview"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type productHandler struct {
	upsertHandler      productview.UpsertProductCommandHandler
	deleteHandler      productview.DeleteProductCommandHandler
	updateImageHandler productview.UpdateImageURLsCommandHandler
}

func newProductHandler(
	upsertHandler productview.UpsertProductCommandHandler,
	deleteHandler productview.DeleteProductCommandHandler,
	updateImageHandler productview.UpdateImageURLsCommandHandler,
) *productHandler {
	return &productHandler{
		upsertHandler:      upsertHandler,
		deleteHandler:      deleteHandler,
		updateImageHandler: updateImageHandler,
	}
}

func (h *productHandler) HandleProductUpdated(ctx context.Context, e *catalog_events.ProductUpdatedEvent) error {
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

	if err := h.upsertHandler.Handle(ctx, productview.UpsertProductCommand{Product: view}); err != nil {
		return err
	}

	h.log(ctx).Debug("product view updated",
		zap.String("productID", e.Payload.ProductID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

// mapAttributes converts event attributes to domain attributes and builds attrs map.
func mapAttributes(eventAttrs *[]catalog_events.AttributeValue) ([]productview.AttributeValue, map[string]any) {
	if eventAttrs == nil || len(*eventAttrs) == 0 {
		return nil, nil
	}
	return lo.Map(*eventAttrs, mapAttribute), buildAttrsMap(*eventAttrs)
}

// mapAttribute converts a single event attribute to a domain attribute.
func mapAttribute(attr catalog_events.AttributeValue, _ int) productview.AttributeValue {
	var optionSlugValues []string
	if attr.OptionSlugValues != nil {
		optionSlugValues = *attr.OptionSlugValues
	}

	return productview.AttributeValue{
		AttributeID:      attr.AttributeID,
		Slug:             attr.AttributeSlug,
		OptionSlugValue:  attr.OptionSlugValue,
		OptionSlugValues: optionSlugValues,
		NumericValue:     attr.NumericValue,
		TextValue:        attr.TextValue,
		BooleanValue:     attr.BooleanValue,
	}
}

// buildAttrsMap builds a denormalized attributes map for filtering (slug -> value).
func buildAttrsMap(eventAttrs []catalog_events.AttributeValue) map[string]any {
	attrs := make(map[string]any, len(eventAttrs))
	for _, attr := range eventAttrs {
		if attr.AttributeSlug == "" {
			continue
		}
		if attr.NumericValue != nil {
			attrs[attr.AttributeSlug] = *attr.NumericValue
		} else if attr.OptionSlugValues != nil && len(*attr.OptionSlugValues) > 0 {
			attrs[attr.AttributeSlug] = *attr.OptionSlugValues
		} else if attr.OptionSlugValue != nil {
			attrs[attr.AttributeSlug] = *attr.OptionSlugValue
		} else if attr.TextValue != nil {
			attrs[attr.AttributeSlug] = *attr.TextValue
		} else if attr.BooleanValue != nil {
			attrs[attr.AttributeSlug] = *attr.BooleanValue
		}
	}
	return attrs
}

func (h *productHandler) HandleProductDeleted(ctx context.Context, e *catalog_events.ProductDeletedEvent) error {
	cmd := productview.DeleteProductCommand{ProductID: e.Payload.ProductID}
	if err := h.deleteHandler.Handle(ctx, cmd); err != nil {
		return err
	}

	h.log(ctx).Debug("product view deleted",
		zap.String("productID", e.Payload.ProductID),
		zap.String("eventID", e.Metadata.EventID))

	return nil
}

func (h *productHandler) HandleProductImagePromoted(ctx context.Context, e *image_events.ProductImagePromotedEvent) error {
	cmd := productview.UpdateImageURLsCommand{
		ProductID:     e.Payload.ProductID,
		ImageID:       e.Payload.ImageID,
		SmallImageURL: e.Payload.SmallImageURL,
		LargeImageURL: e.Payload.LargeImageURL,
	}
	if err := h.updateImageHandler.Handle(ctx, cmd); err != nil {
		return err
	}

	h.log(ctx).Debug("product image URLs updated",
		zap.String("productID", e.Payload.ProductID),
		zap.String("imageID", e.Payload.ImageID),
		zap.String("eventID", e.Metadata.EventID))

	return nil
}

func (h *productHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "product-handler"))
}
