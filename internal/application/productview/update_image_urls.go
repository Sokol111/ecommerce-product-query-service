package productview

import (
	"context"
	"fmt"
)

type UpdateImageURLsCommand struct {
	ProductID     string
	ImageID       string
	SmallImageURL string
	LargeImageURL string
}

type UpdateImageURLsCommandHandler interface {
	Handle(ctx context.Context, cmd UpdateImageURLsCommand) error
}

type updateImageURLsHandler struct {
	repo Repository
}

func NewUpdateImageURLsHandler(repo Repository) UpdateImageURLsCommandHandler {
	return &updateImageURLsHandler{repo: repo}
}

func (h *updateImageURLsHandler) Handle(ctx context.Context, cmd UpdateImageURLsCommand) error {
	if err := h.repo.UpdateImageURLs(ctx, cmd.ProductID, cmd.ImageID, cmd.SmallImageURL, cmd.LargeImageURL); err != nil {
		return fmt.Errorf("failed to update product image URLs: %w", err)
	}
	return nil
}
