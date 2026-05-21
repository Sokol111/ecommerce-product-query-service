package attributeview

import (
	"context"
	"fmt"
)

type UpsertAttributeCommand struct {
	Attribute *AttributeView
}

type UpsertAttributeCommandHandler interface {
	Handle(ctx context.Context, cmd UpsertAttributeCommand) error
}

type upsertAttributeHandler struct {
	repo Repository
}

func NewUpsertAttributeHandler(repo Repository) UpsertAttributeCommandHandler {
	return &upsertAttributeHandler{repo: repo}
}

func (h *upsertAttributeHandler) Handle(ctx context.Context, cmd UpsertAttributeCommand) error {
	if err := h.repo.Upsert(ctx, cmd.Attribute); err != nil {
		return fmt.Errorf("failed to upsert attribute view: %w", err)
	}
	return nil
}
