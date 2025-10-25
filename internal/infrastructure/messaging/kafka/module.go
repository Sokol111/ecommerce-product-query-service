package kafka

import (
	"github.com/Sokol111/ecommerce-commons/pkg/messaging"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		consumer.RegisterHandlerAndConsumer[messaging.ProductCreated]("productCreatedHandler", newProductCreatedHandler),
		consumer.RegisterHandlerAndConsumer[messaging.ProductUpdated]("productUpdatedHandler", newProductUpdatedHandler),
	)
}
