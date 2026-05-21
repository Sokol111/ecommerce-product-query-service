package productview

import (
	"context"
	"fmt"
)

type DeleteProductCommand struct {
	ProductID string
}

type DeleteProductCommandHandler interface {
	Handle(ctx context.Context, cmd DeleteProductCommand) error
}

type deleteProductHandler struct {
	repo Repository
}

func NewDeleteProductHandler(repo Repository) DeleteProductCommandHandler {
	return &deleteProductHandler{repo: repo}
}

func (h *deleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	if err := h.repo.Delete(ctx, cmd.ProductID); err != nil {
		return fmt.Errorf("failed to delete product view: %w", err)
	}
	return nil
}
