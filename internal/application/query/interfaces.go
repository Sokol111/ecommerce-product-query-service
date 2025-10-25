package query

import (
	"context"

	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type GetProductByIDQuery struct {
	ID string
}

type GetRandomProductsQuery struct {
	Amount int
}

type GetProductByIDQueryHandler interface {
	Handle(ctx context.Context, query GetProductByIDQuery) (*productview.ProductView, error)
}

type GetRandomProductsQueryHandler interface {
	Handle(ctx context.Context, query GetRandomProductsQuery) ([]*productview.ProductView, error)
}
