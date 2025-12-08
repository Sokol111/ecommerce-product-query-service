package http

import (
	"context"
	"errors"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	"github.com/Sokol111/ecommerce-product-query-service-api/api"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/query"
)

type productHandler struct {
	getByIDHandler   query.GetProductByIDQueryHandler
	getRandomHandler query.GetRandomProductsQueryHandler
}

func newProductHandler(
	getByIDHandler query.GetProductByIDQueryHandler,
	getRandomHandler query.GetRandomProductsQueryHandler,
) api.StrictServerInterface {
	return &productHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
	}
}

func (h *productHandler) GetProductById(c context.Context, request api.GetProductByIdRequestObject) (api.GetProductByIdResponseObject, error) {
	q := query.GetProductByIDQuery{ID: request.Id}

	found, err := h.getByIDHandler.Handle(c, q)
	if errors.Is(err, persistence.ErrEntityNotFound) {
		return api.GetProductById404ApplicationProblemPlusJSONResponse{
			Status: 404,
			Type:   "about:blank",
			Title:  "Product not found",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return api.GetProductById200JSONResponse{
		Id:       found.ID,
		Name:     found.Name,
		Price:    found.Price,
		Quantity: found.Quantity,
	}, nil
}

func (h *productHandler) GetRandomProducts(c context.Context, request api.GetRandomProductsRequestObject) (api.GetRandomProductsResponseObject, error) {
	q := query.GetRandomProductsQuery{Count: request.Params.Count}

	products, err := h.getRandomHandler.Handle(c, q)
	if err != nil {
		return nil, err
	}

	response := make([]api.ProductResponse, len(products))
	for i, p := range products {
		response[i] = api.ProductResponse{
			Id:       p.ID,
			Name:     p.Name,
			Price:    p.Price,
			Quantity: p.Quantity,
		}
	}

	return api.GetRandomProducts200JSONResponse(response), nil
}
