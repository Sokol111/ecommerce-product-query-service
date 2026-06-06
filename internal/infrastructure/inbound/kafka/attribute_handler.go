package kafka

import (
	"context"

	"github.com/samber/lo"

	catalog_events "github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/attributeview"
	"go.uber.org/zap"
)

type attributeHandler struct {
	upsertHandler attributeview.UpsertAttributeCommandHandler
}

func newAttributeHandler(upsertHandler attributeview.UpsertAttributeCommandHandler) *attributeHandler {
	return &attributeHandler{
		upsertHandler: upsertHandler,
	}
}

func (h *attributeHandler) HandleAttributeUpdated(ctx context.Context, e *catalog_events.AttributeUpdatedEvent) error {
	view := attributeview.Reconstruct(
		e.Payload.AttributeID,
		e.Payload.Version,
		e.Payload.Slug,
		e.Payload.Name,
		attributeview.AttributeType(e.Payload.Type),
		e.Payload.Unit,
		e.Payload.Enabled,
		e.Payload.ModifiedAt,
		lo.Map(e.Payload.Options, mapOption),
	)

	cmd := attributeview.UpsertAttributeCommand{Attribute: view}
	if err := h.upsertHandler.Handle(ctx, cmd); err != nil {
		return err
	}

	h.log(ctx).Debug("attribute view updated",
		zap.String("attributeID", e.Payload.AttributeID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func mapOption(opt catalog_events.AttributeOption, _ int) attributeview.AttributeOption {
	return attributeview.AttributeOption{
		Slug:      opt.Slug,
		Name:      opt.Name,
		ColorCode: opt.ColorCode,
		SortOrder: opt.SortOrder,
	}
}

func (h *attributeHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx)
}
