package main

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/commonsconfig"
	"github.com/Sokol111/ecommerce-commons/pkg/commonsgin"
	"github.com/Sokol111/ecommerce-commons/pkg/commonslogging"
	"github.com/Sokol111/ecommerce-commons/pkg/commonsmongo"
	"github.com/Sokol111/ecommerce-commons/pkg/commonsserver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	commonslogging.ZapLoggingModule,
	commonsconfig.ViperModule,
	commonsmongo.MongoModule,
	commonsgin.GinModule,
	commonsserver.HttpServerModule,
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
