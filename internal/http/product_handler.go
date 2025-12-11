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
}

func newProductHandler(
	getByIDHandler query.GetProductByIDQueryHandler,
	getRandomHandler query.GetRandomProductsQueryHandler,
) httpapi.StrictServerInterface {
	return &productHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
	}
}

func toProductResponse(p *productview.ProductView) httpapi.ProductResponse {
	return httpapi.ProductResponse{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		ImageId:     p.ImageID,
		CategoryId:  p.CategoryID,
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
