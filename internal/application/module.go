package application

import (
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/productview"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
func Module() fx.Option {
	return fx.Options(
		// Query handlers
		fx.Provide(
			productview.NewGetProductByIDHandler,
			productview.NewGetRandomProductsHandler,
			productview.NewGetListProductsHandler,
			productview.NewGetProductFacetsHandler,
		),
		// Command handlers
		fx.Provide(
			productview.NewUpsertProductHandler,
			productview.NewDeleteProductHandler,
			productview.NewUpdateImageURLsHandler,
			attributeview.NewUpsertAttributeHandler,
		),
	)
}
