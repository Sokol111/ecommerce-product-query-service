package main

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/module"
	"github.com/Sokol111/ecommerce-commons/pkg/swaggerui"
	"github.com/Sokol111/ecommerce-product-query-service-api/api"
	"github.com/Sokol111/ecommerce-product-query-service/internal/http"
	"github.com/Sokol111/ecommerce-product-query-service/internal/kafka"
	"github.com/Sokol111/ecommerce-product-query-service/internal/productdetail"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	module.NewInfraModule(),
	module.NewKafkaModule(),
	kafka.NewKafkaHandlerModule(),
	http.NewHttpHandlerModule(),
	productdetail.NewCategoryListViewModule(),
	swaggerui.NewSwaggerModule(swaggerui.SwaggerConfig{OpenAPIContent: api.OpenAPIDoc}),
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
