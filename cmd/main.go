package main

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/config"
	"github.com/Sokol111/ecommerce-commons/pkg/gin"
	commonskafka "github.com/Sokol111/ecommerce-commons/pkg/kafka"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/outbox"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/producer"
	"github.com/Sokol111/ecommerce-commons/pkg/logging"
	"github.com/Sokol111/ecommerce-commons/pkg/mongo"
	"github.com/Sokol111/ecommerce-commons/pkg/server"
	"github.com/Sokol111/ecommerce-product-query-service/internal/kafka"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	logging.ZapLoggingModule,
	config.ViperModule,
	mongo.MongoModule,
	gin.GinModule,
	server.HttpServerModule,
	commonskafka.KafkaModule,
	producer.ProducerModule,
	outbox.OutboxModule,
	consumer.ConsumerModule,
	kafka.HandlersModule,
)

func main() {
	app := fx.New(
		AppModules,
		fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("Application starting...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("Application stopping...")
					return nil
				},
			})
		}),
	)
	app.Run()
}
