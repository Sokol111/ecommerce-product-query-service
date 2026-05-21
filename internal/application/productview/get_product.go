package productview

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

type GetProductByIDQuery struct {
	ID string
}

type GetProductByIDQueryHandler interface {
	Handle(ctx context.Context, query GetProductByIDQuery) (*ProductView, error)
}

type getProductByIDHandler struct {
	repo Repository
}

func NewGetProductByIDHandler(repo Repository) GetProductByIDQueryHandler {
	return &getProductByIDHandler{repo: repo}
}

func (h *getProductByIDHandler) Handle(ctx context.Context, query GetProductByIDQuery) (*ProductView, error) {
	p, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrEntityNotFound) {
			return nil, mongo.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to get product view: %w", err)
	}
	return p, nil
}
