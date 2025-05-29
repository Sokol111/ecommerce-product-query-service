package handler

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type productCreatedHandler struct {
	log *zap.Logger
}

// func ProvideProductCreatedConsumer(
// 	lc fx.Lifecycle,
// 	log *zap.Logger,
// 	conf config.Config,
// 	handler consumer.Handler[payload.ProductCreated],
// ) (consumer.Consumer, error) {
// 	c, err := consumer.ProvideNewConsumer(lc, log, conf, handler, "productCreatedHandler")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return consumer.Consumer(c), nil
// }

func NewProductCreatedHandler(log *zap.Logger) consumer.Handler[payload.ProductCreated] {
	return &productCreatedHandler{
		log: log,
	}
}

func (h *productCreatedHandler) Process(ctx context.Context, e *event.Event[payload.ProductCreated]) error {
	h.log.Info("message received", zap.String("message", fmt.Sprintf("%v", e)))
	return nil
}

func (h *productCreatedHandler) Validate(payload *payload.ProductCreated) error {
	return nil
}
