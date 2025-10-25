package query

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type getRandomProductsHandler struct {
	repo productview.Repository
}

func NewGetRandomProductsHandler(repo productview.Repository) GetRandomProductsQueryHandler {
	return &getRandomProductsHandler{repo: repo}
}

func (h *getRandomProductsHandler) Handle(ctx context.Context, query GetRandomProductsQuery) ([]*productview.ProductView, error) {
	products, err := h.repo.FindRandom(ctx, query.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to get random products: %w", err)
	}
	return products, nil
}
