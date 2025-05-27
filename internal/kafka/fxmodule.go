package kafka

import (
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/kafka/handler"
	"go.uber.org/fx"
)

var HandlersModule = fx.Options(
	productCreatedHandlerModule,
)

var productCreatedHandlerModule = fx.Module(
	"productCreatedHandlerModule",
	fx.Provide(
		handler.NewProductCreatedHandler,
		fx.Private,
	),
	fx.Provide(
		fx.Annotate(
			func(handler consumer.Handler[any]) consumer.HandlerDef {
				return consumer.HandlerDef{Name: "productCreatedHandler", Handler: handler}
			},
			fx.ResultTags(`group:"kafka_handlers"`)),
	),
)
