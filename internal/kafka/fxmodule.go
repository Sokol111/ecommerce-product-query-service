package kafka

import (
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/kafka/handler"
	"go.uber.org/fx"
)

var HandlersModule = fx.Options(
	consumer.RegisterHandlerAndConsumer[payload.ProductCreated]("productCreatedHandler", handler.NewProductCreatedHandler),
	consumer.RegisterHandlerAndConsumer[payload.ProductUpdated]("productUpdatedHandler", handler.NewProductUpdatedHandler),
)
