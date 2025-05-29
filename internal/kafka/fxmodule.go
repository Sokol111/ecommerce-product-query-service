package kafka

import (
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"github.com/Sokol111/ecommerce-product-query-service/internal/kafka/handler"
	"go.uber.org/fx"
)

// var HandlersModule = fx.Options(
// 	fx.Provide(
// 		handler.NewProductCreatedHandler,
// 	),
// 	fx.Provide(
// 		fx.Annotate(
// 			handler.ProvideProductCreatedConsumer,
// 			fx.ResultTags(`group:"consumers"`),
// 		),
// 	),
// 	fx.Invoke(
// 		fx.Annotate(
// 			func(consumers []consumer.Consumer, log *zap.Logger) {
// 				log.Info("Kafka consumers initialized", zap.Int("len", len(consumers)))
// 			},
// 			fx.ParamTags(`group:"consumers"`),
// 		),
// 	),
// )

var HandlersModule = fx.Options(
	consumer.RegisterHandlerAndConsumer[payload.ProductCreated]("productCreatedHandler", handler.NewProductCreatedHandler),
)
