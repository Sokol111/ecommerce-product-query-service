package productview

import (
	"context"
	"fmt"
)

type GetProductFacetsQuery struct {
	CategoryID string
}

type GetProductFacetsQueryHandler interface {
	Handle(ctx context.Context, query GetProductFacetsQuery) (*FacetsResult, error)
}

type getProductFacetsHandler struct {
	repo Repository
}

func NewGetProductFacetsHandler(repo Repository) GetProductFacetsQueryHandler {
	return &getProductFacetsHandler{repo: repo}
}

func (h *getProductFacetsHandler) Handle(ctx context.Context, query GetProductFacetsQuery) (*FacetsResult, error) {
	result, err := h.repo.FindFacets(ctx, query.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product facets: %w", err)
	}

	return result, nil
}
