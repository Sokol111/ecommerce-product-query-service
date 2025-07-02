package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/model"
)

type productUpdatedHandler struct {
	productDetailService model.ProductDetailService
}

func newProductUpdatedHandler(productDetailService model.ProductDetailService) consumer.Handler[payload.ProductUpdated] {
	return &productUpdatedHandler{
		productDetailService: productDetailService,
	}
}

func (h *productUpdatedHandler) Process(ctx context.Context, e *event.Event[payload.ProductUpdated]) error {
	return h.productDetailService.ProcessProductUpdatedEvent(ctx, e)
}

func (h *productUpdatedHandler) Validate(payload *payload.ProductUpdated) error {
	return nil
}
