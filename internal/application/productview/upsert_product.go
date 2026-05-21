package productview

import (
	"context"
	"fmt"
)

type UpsertProductCommand struct {
	Product *ProductView
}

type UpsertProductCommandHandler interface {
	Handle(ctx context.Context, cmd UpsertProductCommand) error
}

type upsertProductHandler struct {
	repo Repository
}

func NewUpsertProductHandler(repo Repository) UpsertProductCommandHandler {
	return &upsertProductHandler{repo: repo}
}

func (h *upsertProductHandler) Handle(ctx context.Context, cmd UpsertProductCommand) error {
	if err := h.repo.Upsert(ctx, cmd.Product); err != nil {
		return fmt.Errorf("failed to upsert product view: %w", err)
	}
	return nil
}
