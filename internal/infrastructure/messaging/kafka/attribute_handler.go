package kafka

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	catalog_events "github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/attributeview"
	"go.uber.org/zap"
)

type attributeHandler struct {
	repo attributeview.Repository
}

func newAttributeHandler(repo attributeview.Repository) *attributeHandler {
	return &attributeHandler{
		repo: repo,
	}
}

func (h *attributeHandler) Process(ctx context.Context, event any) error {
	switch evt := event.(type) {
	case *catalog_events.AttributeUpdatedEvent:
		return h.handleAttributeUpdated(ctx, evt)

	default:
		logger.Get(ctx).Warn("unknown event type, skipping",
			zap.String("type", fmt.Sprintf("%T", event)))
		return fmt.Errorf("unhandled event type: %T: %w", event, consumer.ErrSkipMessage)
	}
}

func (h *attributeHandler) handleAttributeUpdated(ctx context.Context, e *catalog_events.AttributeUpdatedEvent) error {
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

	if err := h.repo.Upsert(ctx, view); err != nil {
		return fmt.Errorf("failed to upsert attribute view: %w", err)
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
