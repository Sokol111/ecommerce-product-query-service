package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/Sokol111/ecommerce-product-query-service-api/api"
	"github.com/Sokol111/ecommerce-product-query-service/internal/model"
)

type productHandler struct {
	service model.ProductDetailService
}

func newProductHandler(service model.ProductDetailService) api.StrictServerInterface {
	return &productHandler{service}
}

func (h *productHandler) GetProductById(c context.Context, request api.GetProductByIdRequestObject) (api.GetProductByIdResponseObject, error) {
	found, err := h.service.GetById(c, request.Id)
	if errors.Is(err, model.ErrNotFound) {
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
