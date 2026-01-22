package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/query"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/attributeview"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type productHandler struct {
	getByIDHandler   query.GetProductByIDQueryHandler
	getRandomHandler query.GetRandomProductsQueryHandler
	getListHandler   query.GetListProductsQueryHandler
	attributeRepo    attributeview.Repository
}

func newProductHandler(
	getByIDHandler query.GetProductByIDQueryHandler,
	getRandomHandler query.GetRandomProductsQueryHandler,
	getListHandler query.GetListProductsQueryHandler,
	attributeRepo attributeview.Repository,
) httpapi.Handler {
	return &productHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
		getListHandler:   getListHandler,
		attributeRepo:    attributeRepo,
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

func toOptBool(b *bool) httpapi.OptBool {
	if b == nil {
		return httpapi.OptBool{}
	}
	return httpapi.NewOptBool(*b)
}

// toProductAttributes joins product attributes with master data and maps to HTTP response
func (h *productHandler) toProductAttributes(ctx context.Context, attrs []productview.ProductAttribute) ([]httpapi.ProductAttribute, error) {
	if len(attrs) == 0 {
		return nil, nil
	}

	attrIDs := lo.Uniq(lo.Map(attrs, func(attr productview.ProductAttribute, _ int) string {
		return attr.AttributeID
	}))

	masterAttrs, err := h.attributeRepo.FindByIDs(ctx, attrIDs)
	if err != nil {
		return nil, err
	}

	attrByID := lo.KeyBy(masterAttrs, func(attr *attributeview.AttributeView) string { return attr.ID })

	return lo.FilterMap(attrs, func(attr productview.ProductAttribute, i int) (httpapi.ProductAttribute, bool) {
		master, ok := attrByID[attr.AttributeID]
		// Skip if master data not found or attribute is disabled
		if !ok || !master.Enabled {
			return httpapi.ProductAttribute{}, false
		}

		result := httpapi.ProductAttribute{
			AttributeId:      attr.AttributeID,
			Slug:             attr.Slug,
			Name:             master.Name,
			Type:             httpapi.ProductAttributeType(master.Type),
			Unit:             toOptString(master.Unit),
			Role:             httpapi.ProductAttributeRole("specification"),
			SortOrder:        i,
			OptionSlugValue:  toOptString(attr.OptionSlugValue),
			OptionSlugValues: attr.OptionSlugValues,
			NumericValue:     toOptFloat64(attr.NumericValue),
			TextValue:        toOptString(attr.TextValue),
			BooleanValue:     toOptBool(attr.BooleanValue),
		}

		optionsBySlug := lo.KeyBy(master.Options, func(opt attributeview.AttributeOption) string { return opt.Slug })

		if attr.OptionSlugValue != nil {
			if opt, ok := optionsBySlug[*attr.OptionSlugValue]; ok {
				result.OptionName = toOptString(&opt.Name)
				result.OptionColorCode = toOptString(opt.ColorCode)
			}
		}

		if len(attr.OptionSlugValues) > 0 {
			result.OptionNames = lo.FilterMap(attr.OptionSlugValues, func(slug string, _ int) (string, bool) {
				opt, ok := optionsBySlug[slug]
				return opt.Name, ok
			})
		}

		return result, true
	}), nil
}

func (h *productHandler) toProductResponse(ctx context.Context, p *productview.ProductView) (*httpapi.ProductResponse, error) {
	attrs, err := h.toProductAttributes(ctx, p.Attributes)
	if err != nil {
		return nil, err
	}

	return &httpapi.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: toOptString(p.Description),
		Price:       float64(p.Price),
		Quantity:    p.Quantity,
		ImageId:     toOptString(p.ImageID),
		ImageUrl:    toOptString(p.ImageURL),
		CategoryId:  toOptString(p.CategoryID),
		Attributes:  attrs,
	}, nil
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

	resp, err := h.toProductResponse(ctx, found)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *productHandler) GetRandomProducts(ctx context.Context, params httpapi.GetRandomProductsParams) (httpapi.GetRandomProductsRes, error) {
	q := query.GetRandomProductsQuery{Count: params.Count}

	products, err := h.getRandomHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	response := make([]httpapi.ProductResponse, len(products))
	for i, p := range products {
		resp, err := h.toProductResponse(ctx, p)
		if err != nil {
			return nil, err
		}
		response[i] = *resp
	}

	return (*httpapi.GetRandomProductsOKApplicationJSON)(&response), nil
}

func (h *productHandler) GetProductList(ctx context.Context, params httpapi.GetProductListParams) (httpapi.GetProductListRes, error) {
	var categoryID *string
	if params.CategoryId.IsSet() {
		categoryID = &params.CategoryId.Value
	}

	// Parse price filters
	var minPrice, maxPrice *float64
	if params.MinPrice.IsSet() {
		minPrice = &params.MinPrice.Value
	}
	if params.MaxPrice.IsSet() {
		maxPrice = &params.MaxPrice.Value
	}

	// Parse attribute filters from JSON string
	var attributeFilters []query.AttributeFilter
	if params.AttributeFilters.IsSet() && params.AttributeFilters.Value != "" {
		if err := json.Unmarshal([]byte(params.AttributeFilters.Value), &attributeFilters); err != nil {
			return &httpapi.GetProductListBadRequest{
				Status: 400,
				Type:   *aboutBlankURL,
				Title:  "Invalid attributeFilters format",
				Detail: httpapi.NewOptString("attributeFilters must be a valid JSON array"),
			}, nil
		}
	}

	q := query.GetListProductsQuery{
		Page:             params.Page,
		Size:             params.Size,
		CategoryID:       categoryID,
		Sort:             string(params.Sort.Or(httpapi.GetProductListSortName)),
		Order:            string(params.Order.Or(httpapi.GetProductListOrderDesc)),
		MinPrice:         minPrice,
		MaxPrice:         maxPrice,
		AttributeFilters: attributeFilters,
	}

	result, err := h.getListHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	items := make([]httpapi.ProductResponse, len(result.Items))
	for i, p := range result.Items {
		resp, err := h.toProductResponse(ctx, p)
		if err != nil {
			return nil, err
		}
		items[i] = *resp
	}

	return &httpapi.ProductListResponse{
		Items: items,
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}, nil
}
