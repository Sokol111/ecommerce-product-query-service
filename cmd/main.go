package main

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/modules"
	"github.com/Sokol111/ecommerce-commons/pkg/swaggerui"
	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application"
	"github.com/Sokol111/ecommerce-product-query-service/internal/http"
	"github.com/Sokol111/ecommerce-product-query-service/internal/infrastructure/client"
	"github.com/Sokol111/ecommerce-product-query-service/internal/infrastructure/messaging/kafka"
	"github.com/Sokol111/ecommerce-product-query-service/internal/infrastructure/persistence/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	// Infrastructure
	modules.NewCoreModule(),
	modules.NewPersistenceModule(),
	modules.NewHTTPModule(),
	modules.NewObservabilityModule(),
	modules.NewMessagingModule(),

	// Domain & Application
	mongo.Module(),
	client.AttributeClientModule(),
	application.Module(),

	// Messaging
	kafka.Module(),

	// HTTP
	http.NewHttpHandlerModule(),
	swaggerui.NewSwaggerModule(swaggerui.SwaggerConfig{OpenAPIContent: httpapi.OpenAPIDoc}),
)

func main() {
	app := fx.New(
		AppModules,
		fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Info("Application stopping...")
					return nil
				},
			})
		}),
	)
	app.Run()
}
