package productview

import (
	"context"
	"fmt"
)

type GetRandomProductsQuery struct {
	Count int
}

type GetRandomProductsQueryHandler interface {
	Handle(ctx context.Context, query GetRandomProductsQuery) ([]*ProductView, error)
}

type getRandomProductsHandler struct {
	repo Repository
}

func NewGetRandomProductsHandler(repo Repository) GetRandomProductsQueryHandler {
	return &getRandomProductsHandler{repo: repo}
}

func (h *getRandomProductsHandler) Handle(ctx context.Context, query GetRandomProductsQuery) ([]*ProductView, error) {
	products, err := h.repo.FindRandom(ctx, query.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to get random products: %w", err)
	}
	return products, nil
}
