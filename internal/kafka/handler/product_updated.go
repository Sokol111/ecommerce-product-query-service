package handler

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type productUpdatedHandler struct {
	log *zap.Logger
}

func newProductUpdatedHandler(log *zap.Logger) consumer.Handler[payload.ProductUpdated] {
	return &productUpdatedHandler{
		log: log,
	}
}

func (h *productUpdatedHandler) Process(ctx context.Context, e *event.Event[payload.ProductUpdated]) error {
	h.log.Info("message received", zap.String("message", fmt.Sprintf("%v", e)))
	return nil
}

func (h *productUpdatedHandler) Validate(payload *payload.ProductUpdated) error {
	return nil
}
