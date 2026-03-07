package query

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type GetProductFacetsQuery struct {
	CategoryID string
}

type GetProductFacetsQueryHandler interface {
	Handle(ctx context.Context, query GetProductFacetsQuery) (*productview.FacetsResult, error)
}

type getProductFacetsHandler struct {
	repo productview.Repository
}

func NewGetProductFacetsHandler(repo productview.Repository) GetProductFacetsQueryHandler {
	return &getProductFacetsHandler{repo: repo}
}

func (h *getProductFacetsHandler) Handle(ctx context.Context, query GetProductFacetsQuery) (*productview.FacetsResult, error) {
	result, err := h.repo.FindFacets(ctx, query.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product facets: %w", err)
	}

	return result, nil
}
