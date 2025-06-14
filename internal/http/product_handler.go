package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/Sokol111/ecommerce-product-query-service-api/api"
	"github.com/Sokol111/ecommerce-product-query-service/pkg/product"
)

type productHandler struct {
	service product.Service
}

func newProductHandler(service product.Service) api.StrictServerInterface {
	return &productHandler{service}
}

func (h *productHandler) GetProductById(c context.Context, request api.GetProductByIdRequestObject) (api.GetProductByIdResponseObject, error) {
	found, err := h.service.GetById(c, request.Id)
	if errors.Is(err, product.NotFoundError) {
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

func (h *productHandler) GetAll(c context.Context, _ api.GetAllRequestObject) (api.GetAllResponseObject, error) {
	products, err := h.service.GetAll(c)
	if err != nil {
		return api.GetAll500JSONResponse{Code: 500, Message: http.StatusText(500)}, err
	}
	response := make(api.GetAll200JSONResponse, 0, len(products))
	for _, u := range products {
		response = append(response, api.ProductResponse{
			Id:       u.ID,
			Name:     u.Name,
			Price:    u.Price,
			Quantity: u.Quantity,
		})
	}
	return response, nil
}
