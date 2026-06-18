package connect

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/Sokol111/ecommerce-commons/pkg/security/validation"
	productqueryv1connect "github.com/Sokol111/ecommerce-product-query-service-api/gen/connect/product_query/v1/productqueryv1connect"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/productview"
	"go.uber.org/fx"
)

// Module provides the Connect gRPC/Connect-RPC server handler for product query operations.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			newProductQueryHandler,
			provideProcedurePermissions,
		),
		fx.Invoke(registerConnectRoutes),
	)
}

func newProductQueryHandler(
	getByIDHandler productview.GetProductByIDQueryHandler,
	getRandomHandler productview.GetRandomProductsQueryHandler,
	getListHandler productview.GetListProductsQueryHandler,
	getFacetsHandler productview.GetProductFacetsQueryHandler,
	attributeRepo attributeview.Repository,
) *productQueryHandler {
	return &productQueryHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
		getListHandler:   getListHandler,
		getFacetsHandler: getFacetsHandler,
		attributeRepo:    attributeRepo,
	}
}

func registerConnectRoutes(
	mux *http.ServeMux,
	handler *productQueryHandler,
	interceptors []connect.Interceptor,
) {
	path, h := productqueryv1connect.NewProductQueryServiceHandler(handler, connect.WithInterceptors(interceptors...))
	mux.Handle(path, h)
}

// provideProcedurePermissions maps each product query RPC to required permissions.
//
// Permission semantics:
//   - nil         → public endpoint, no authentication required
//   - []string{}  → authenticated, no specific permissions required
//   - []string{…} → authenticated and must hold at least one of the listed permissions
func provideProcedurePermissions() validation.ProcedurePermissions {
	return validation.ProcedurePermissions{
		productqueryv1connect.ProductQueryServiceGetProductByIdProcedure:    nil,
		productqueryv1connect.ProductQueryServiceGetRandomProductsProcedure: nil,
		productqueryv1connect.ProductQueryServiceGetProductListProcedure:    nil,
		productqueryv1connect.ProductQueryServiceGetProductFacetsProcedure:  nil,
	}
}
