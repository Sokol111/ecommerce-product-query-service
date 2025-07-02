package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/model"
)

type productCreatedHandler struct {
	productDetailService model.ProductDetailService
}

func newProductCreatedHandler(productDetailService model.ProductDetailService) consumer.Handler[payload.ProductCreated] {
	return &productCreatedHandler{
		productDetailService: productDetailService,
	}
}

func (h *productCreatedHandler) Process(ctx context.Context, e *event.Event[payload.ProductCreated]) error {
	return h.productDetailService.ProcessProductCreatedEvent(ctx, e)
}

func (h *productCreatedHandler) Validate(payload *payload.ProductCreated) error {
	return nil
}
