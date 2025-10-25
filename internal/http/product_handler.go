package http

import (
	"context"
	"errors"
	"net/http"

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
		return api.GetProductById404JSONResponse{Code: 404, Message: "Product not found"}, nil
	}
	if err != nil {
		return api.GetProductById500JSONResponse{Code: 500, Message: http.StatusText(500)}, err
	}

	return api.GetProductById200JSONResponse{
		Id:       found.ID,
		Name:     found.Name,
		Price:    found.Price,
		Quantity: found.Quantity,
	}, nil
}
