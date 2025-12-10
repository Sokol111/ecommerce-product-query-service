package http

import (
	"fmt"

	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHttpHandlerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			newProductHandler,
			func(ssi httpapi.StrictServerInterface) httpapi.ServerInterface {
				return httpapi.NewStrictHandler(ssi, nil)
			},
			// Provide OpenAPI spec - gin module will auto-register validation middleware
			newOpenAPISpec,
		),
		fx.Invoke(registerRoutes),
	)
}

func newOpenAPISpec(log *zap.Logger) (*openapi3.T, error) {
	swagger, err := httpapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}
	return swagger, nil
}

func registerRoutes(engine *gin.Engine, serverInterface httpapi.ServerInterface) {
	httpapi.RegisterHandlers(engine, serverInterface)
}
