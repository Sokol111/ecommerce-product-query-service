package query

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type GetListProductsQuery struct {
	Page       int
	Size       int
	CategoryID *string
	Sort       string
	Order      string
}

type ListProductsResult struct {
	Items []*productview.ProductView
	Page  int
	Size  int
	Total int64
}

type GetListProductsQueryHandler interface {
	Handle(ctx context.Context, query GetListProductsQuery) (*ListProductsResult, error)
}

type getListProductsHandler struct {
	repo productview.Repository
}

func NewGetListProductsHandler(repo productview.Repository) GetListProductsQueryHandler {
	return &getListProductsHandler{repo: repo}
}

func (h *getListProductsHandler) Handle(ctx context.Context, query GetListProductsQuery) (*ListProductsResult, error) {
	listQuery := productview.ListQuery{
		Page:       query.Page,
		Size:       query.Size,
		CategoryID: query.CategoryID,
		Sort:       query.Sort,
		Order:      query.Order,
	}

	result, err := h.repo.FindList(ctx, listQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get products list: %w", err)
	}

	return &ListProductsResult{
		Items: result.Items,
		Page:  result.Page,
		Size:  result.Size,
		Total: result.Total,
	}, nil
}
