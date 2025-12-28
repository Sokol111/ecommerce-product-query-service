package http

import (
	"context"
	"errors"
	"net/url"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/query"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type productHandler struct {
	getByIDHandler   query.GetProductByIDQueryHandler
	getRandomHandler query.GetRandomProductsQueryHandler
	getListHandler   query.GetListProductsQueryHandler
}

func newProductHandler(
	getByIDHandler query.GetProductByIDQueryHandler,
	getRandomHandler query.GetRandomProductsQueryHandler,
	getListHandler query.GetListProductsQueryHandler,
) httpapi.Handler {
	return &productHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
		getListHandler:   getListHandler,
	}
}

var aboutBlankURL, _ = url.Parse("about:blank")

func toOptString(s *string) httpapi.OptString {
	if s == nil {
		return httpapi.OptString{}
	}
	return httpapi.NewOptString(*s)
}

func toOptFloat64(f *float32) httpapi.OptFloat64 {
	if f == nil {
		return httpapi.OptFloat64{}
	}
	return httpapi.NewOptFloat64(float64(*f))
}

func toProductAttributeResponse(attr productview.ProductAttribute, _ int) httpapi.ProductAttribute {
	return httpapi.ProductAttribute{
		AttributeId:  attr.AttributeID,
		Value:        toOptString(attr.Value),
		Values:       attr.Values,
		NumericValue: toOptFloat64(attr.NumericValue),
	}
}

func toProductResponse(p *productview.ProductView) *httpapi.ProductResponse {
	return &httpapi.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: toOptString(p.Description),
		Price:       float64(p.Price),
		Quantity:    p.Quantity,
		ImageId:     toOptString(p.ImageID),
		ImageUrl:    toOptString(p.ImageURL),
		CategoryId:  toOptString(p.CategoryID),
		Attributes:  lo.Map(p.Attributes, toProductAttributeResponse),
	}
}

func (h *productHandler) GetProductById(ctx context.Context, params httpapi.GetProductByIdParams) (httpapi.GetProductByIdRes, error) {
	q := query.GetProductByIDQuery{ID: params.ID}

	found, err := h.getByIDHandler.Handle(ctx, q)
	if errors.Is(err, persistence.ErrEntityNotFound) {
		return &httpapi.GetProductByIdNotFound{
			Status: 404,
			Type:   *aboutBlankURL,
			Title:  "Product not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return toProductResponse(found), nil
}

func (h *productHandler) GetRandomProducts(ctx context.Context, params httpapi.GetRandomProductsParams) (httpapi.GetRandomProductsRes, error) {
	q := query.GetRandomProductsQuery{Count: params.Count}

	products, err := h.getRandomHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	response := lo.Map(products, func(p *productview.ProductView, _ int) httpapi.ProductResponse {
		return *toProductResponse(p)
	})

	return (*httpapi.GetRandomProductsOKApplicationJSON)(&response), nil
}

func (h *productHandler) GetProductList(ctx context.Context, params httpapi.GetProductListParams) (httpapi.GetProductListRes, error) {
	var categoryID *string
	if params.CategoryId.IsSet() {
		categoryID = &params.CategoryId.Value
	}

	q := query.GetListProductsQuery{
		Page:       params.Page,
		Size:       params.Size,
		CategoryID: categoryID,
		Sort:       string(params.Sort.Or(httpapi.GetProductListSortName)),
		Order:      string(params.Order.Or(httpapi.GetProductListOrderDesc)),
	}

	result, err := h.getListHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	return &httpapi.ProductListResponse{
		Items: lo.Map(result.Items, func(p *productview.ProductView, _ int) httpapi.ProductResponse {
			return *toProductResponse(p)
		}),
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}, nil
}
