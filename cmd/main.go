package main

import (
	"context"

	commons_core "github.com/Sokol111/ecommerce-commons/pkg/core"
	commons_http "github.com/Sokol111/ecommerce-commons/pkg/http"
	commons_messaging "github.com/Sokol111/ecommerce-commons/pkg/messaging"
	commons_observability "github.com/Sokol111/ecommerce-commons/pkg/observability"
	commons_persistence "github.com/Sokol111/ecommerce-commons/pkg/persistence"
	commons_pprof "github.com/Sokol111/ecommerce-commons/pkg/pprof"
	commons_swaggerui "github.com/Sokol111/ecommerce-commons/pkg/swaggerui"
	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application"
	"github.com/Sokol111/ecommerce-product-query-service/internal/http"
	"github.com/Sokol111/ecommerce-product-query-service/internal/infrastructure/messaging/kafka"
	"github.com/Sokol111/ecommerce-product-query-service/internal/infrastructure/persistence/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	// Commons
	commons_core.NewCoreModule(),
	commons_persistence.NewPersistenceModule(),
	commons_http.NewHTTPModule(),
	commons_observability.NewObservabilityModule(),
	commons_messaging.NewMessagingModule(),
	commons_pprof.NewPprofModule(),
	commons_swaggerui.NewSwaggerModule(commons_swaggerui.SwaggerConfig{OpenAPIContent: httpapi.OpenAPIDoc}),

	// Application
	mongo.Module(),
	application.Module(),
	kafka.Module(),
	http.Module(),
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
