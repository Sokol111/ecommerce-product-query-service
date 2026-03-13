package http //nolint:revive // intentional package name

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			newProductHandler,
			httpapi.ProvideServer,
		),
		fx.Invoke(registerOgenRoutes),
	)
}

func registerOgenRoutes(mux *http.ServeMux, server *httpapi.Server) {
	mux.Handle("/", server)
}
