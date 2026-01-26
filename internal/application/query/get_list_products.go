package query

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

// AttributeFilter represents a filter for a single attribute
type AttributeFilter struct {
	Slug   string
	Values []string // For single/multiple type attributes
	Min    *float64 // For range type attributes
	Max    *float64 // For range type attributes
}

type GetListProductsQuery struct {
	Page             int
	Size             int
	CategoryID       *string
	Sort             string
	Order            string
	MinPrice         *float64
	MaxPrice         *float64
	AttributeFilters []AttributeFilter
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
		Page:             query.Page,
		Size:             query.Size,
		CategoryID:       query.CategoryID,
		Sort:             query.Sort,
		Order:            query.Order,
		MinPrice:         query.MinPrice,
		MaxPrice:         query.MaxPrice,
		AttributeFilters: lo.Map(query.AttributeFilters, mapAttributeFilter),
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

func mapAttributeFilter(f AttributeFilter, _ int) productview.AttributeFilter {
	return productview.AttributeFilter{
		Slug:   f.Slug,
		Values: f.Values,
		Min:    f.Min,
		Max:    f.Max,
	}
}
