package query

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type GetProductByIDQuery struct {
	ID string
}

type GetProductByIDQueryHandler interface {
	Handle(ctx context.Context, query GetProductByIDQuery) (*productview.ProductView, error)
}

type getProductByIDHandler struct {
	repo productview.Repository
}

func NewGetProductByIDHandler(repo productview.Repository) GetProductByIDQueryHandler {
	return &getProductByIDHandler{repo: repo}
}

func (h *getProductByIDHandler) Handle(ctx context.Context, query GetProductByIDQuery) (*productview.ProductView, error) {
	p, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return nil, persistence.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to get product view: %w", err)
	}
	return p, nil
}
