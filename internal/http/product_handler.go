package http

import (
	"context"
	"errors"

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
) httpapi.StrictServerInterface {
	return &productHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
		getListHandler:   getListHandler,
	}
}

func toProductResponse(p *productview.ProductView) httpapi.ProductResponse {
	var attributes *[]httpapi.ProductAttribute
	if len(p.Attributes) > 0 {
		attrs := make([]httpapi.ProductAttribute, len(p.Attributes))
		for i, attr := range p.Attributes {
			attrs[i] = httpapi.ProductAttribute{
				AttributeId:  attr.AttributeID,
				Value:        attr.Value,
				Values:       &attr.Values,
				NumericValue: attr.NumericValue,
			}
		}
		attributes = &attrs
	}

	return httpapi.ProductResponse{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		ImageId:     p.ImageID,
		ImageUrl:    p.ImageURL,
		CategoryId:  p.CategoryID,
		Attributes:  attributes,
	}
}

func (h *productHandler) GetProductById(c context.Context, request httpapi.GetProductByIdRequestObject) (httpapi.GetProductByIdResponseObject, error) {
	q := query.GetProductByIDQuery{ID: request.Id}

	found, err := h.getByIDHandler.Handle(c, q)
	if errors.Is(err, persistence.ErrEntityNotFound) {
		return httpapi.GetProductById404ApplicationProblemPlusJSONResponse{
			Status: 404,
			Type:   "about:blank",
			Title:  "Product not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return httpapi.GetProductById200JSONResponse(toProductResponse(found)), nil
}

func (h *productHandler) GetRandomProducts(c context.Context, request httpapi.GetRandomProductsRequestObject) (httpapi.GetRandomProductsResponseObject, error) {
	q := query.GetRandomProductsQuery{Count: request.Params.Count}

	products, err := h.getRandomHandler.Handle(c, q)
	if err != nil {
		return nil, err
	}

	response := make([]httpapi.ProductResponse, len(products))
	for i, p := range products {
		response[i] = toProductResponse(p)
	}

	return httpapi.GetRandomProducts200JSONResponse(response), nil
}

func (h *productHandler) GetProductList(c context.Context, request httpapi.GetProductListRequestObject) (httpapi.GetProductListResponseObject, error) {
	// Default sort and order
	sort := "name"
	order := "desc"

	// Override with request params if provided
	if request.Params.Sort != nil {
		sort = string(*request.Params.Sort)
	}

	if request.Params.Order != nil {
		order = string(*request.Params.Order)
	}

	q := query.GetListProductsQuery{
		Page:       request.Params.Page,
		Size:       request.Params.Size,
		CategoryID: request.Params.CategoryId,
		Sort:       sort,
		Order:      order,
	}

	result, err := h.getListHandler.Handle(c, q)
	if err != nil {
		return nil, err
	}

	response := httpapi.GetProductList200JSONResponse{
		Items: make([]httpapi.ProductResponse, 0, len(result.Items)),
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}

	for _, p := range result.Items {
		response.Items = append(response.Items, toProductResponse(p))
	}

	return response, nil
}
