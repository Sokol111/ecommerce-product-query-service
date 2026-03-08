package http //nolint:revive // intentional package name to group HTTP handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/query"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/attributeview"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
)

type productHandler struct {
	getByIDHandler   query.GetProductByIDQueryHandler
	getRandomHandler query.GetRandomProductsQueryHandler
	getListHandler   query.GetListProductsQueryHandler
	getFacetsHandler query.GetProductFacetsQueryHandler
	attributeRepo    attributeview.Repository
}

var aboutBlankURL *url.URL

func init() {
	var err error
	aboutBlankURL, err = url.Parse("about:blank")
	if err != nil {
		panic(fmt.Sprintf("failed to parse about:blank URL: %v", err))
	}
}

func newProductHandler(
	getByIDHandler query.GetProductByIDQueryHandler,
	getRandomHandler query.GetRandomProductsQueryHandler,
	getListHandler query.GetListProductsQueryHandler,
	getFacetsHandler query.GetProductFacetsQueryHandler,
	attributeRepo attributeview.Repository,
) httpapi.Handler {
	return &productHandler{
		getByIDHandler:   getByIDHandler,
		getRandomHandler: getRandomHandler,
		getListHandler:   getListHandler,
		getFacetsHandler: getFacetsHandler,
		attributeRepo:    attributeRepo,
	}
}

func toOptString(s *string) httpapi.OptString {
	if s == nil {
		return httpapi.OptString{}
	}
	return httpapi.NewOptString(*s)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func formatBool(b bool) string {
	return strconv.FormatBool(b)
}

// mapAttributeValue converts a single product attribute to HTTP response using pre-fetched master data
func mapAttributeValue(attr productview.AttributeValue, index int, attrByID map[string]*attributeview.AttributeView) (httpapi.AttributeValue, bool) {
	master, ok := attrByID[attr.AttributeID]
	// Skip if master data not found or attribute is disabled
	if !ok || !master.Enabled {
		return httpapi.AttributeValue{}, false
	}

	optionsBySlug := lo.KeyBy(master.Options, func(opt attributeview.AttributeOption) string { return opt.Slug })

	// Build unified values array based on attribute type
	var values []httpapi.AttributeValueItem

	switch master.Type {
	case "single":
		if attr.OptionSlugValue != nil {
			if opt, ok := optionsBySlug[*attr.OptionSlugValue]; ok {
				values = []httpapi.AttributeValueItem{{
					Slug:      toOptString(attr.OptionSlugValue),
					Value:     opt.Name,
					ColorCode: toOptString(opt.ColorCode),
				}}
			}
		}
	case "multiple":
		values = lo.FilterMap(attr.OptionSlugValues, func(slug string, _ int) (httpapi.AttributeValueItem, bool) {
			opt, ok := optionsBySlug[slug]
			if !ok {
				return httpapi.AttributeValueItem{}, false
			}
			return httpapi.AttributeValueItem{
				Slug:      httpapi.NewOptString(slug),
				Value:     opt.Name,
				ColorCode: toOptString(opt.ColorCode),
			}, true
		})
	case "range":
		if attr.NumericValue != nil {
			values = []httpapi.AttributeValueItem{{
				Value: formatFloat(*attr.NumericValue),
			}}
		}
	case "boolean":
		if attr.BooleanValue != nil {
			values = []httpapi.AttributeValueItem{{
				Value: formatBool(*attr.BooleanValue),
			}}
		}
	case "text":
		if attr.TextValue != nil {
			values = []httpapi.AttributeValueItem{{
				Value: *attr.TextValue,
			}}
		}
	}

	return httpapi.AttributeValue{
		AttributeId: attr.AttributeID,
		Slug:        attr.Slug,
		Name:        master.Name,
		Type:        httpapi.AttributeValueType(master.Type),
		Unit:        toOptString(master.Unit),
		Role:        httpapi.AttributeValueRole("specification"),
		SortOrder:   index,
		Values:      values,
	}, true
}

// toAttributeValuesWithMasterData maps product attributes using pre-fetched master data (no DB call)
func toAttributeValuesWithMasterData(attrs []productview.AttributeValue, attrByID map[string]*attributeview.AttributeView) []httpapi.AttributeValue {
	if len(attrs) == 0 {
		return nil
	}

	return lo.FilterMap(attrs, func(attr productview.AttributeValue, i int) (httpapi.AttributeValue, bool) {
		return mapAttributeValue(attr, i, attrByID)
	})
}

// toProductResponseWithMasterData converts a product to HTTP response using pre-fetched master data (no DB call)
func toProductResponseWithMasterData(p *productview.ProductView, attrByID map[string]*attributeview.AttributeView) *httpapi.ProductResponse {
	return &httpapi.ProductResponse{
		ID:            p.ID,
		Name:          p.Name,
		Description:   toOptString(p.Description),
		Price:         p.Price,
		Quantity:      p.Quantity,
		ImageId:       toOptString(p.ImageID),
		SmallImageUrl: toOptString(p.SmallImageURL),
		LargeImageUrl: toOptString(p.LargeImageURL),
		CategoryId:    toOptString(p.CategoryID),
		Attributes:    toAttributeValuesWithMasterData(p.Attributes, attrByID),
	}
}

// fetchMasterAttributes fetches all unique attributes for a list of products in a single query
func (h *productHandler) fetchMasterAttributes(ctx context.Context, products []*productview.ProductView) (map[string]*attributeview.AttributeView, error) {
	// Collect all unique attribute IDs from all products
	attrIDs := lo.Uniq(lo.FlatMap(products, func(p *productview.ProductView, _ int) []string {
		return lo.Map(p.Attributes, func(attr productview.AttributeValue, _ int) string {
			return attr.AttributeID
		})
	}))

	if len(attrIDs) == 0 {
		return make(map[string]*attributeview.AttributeView), nil
	}

	// Single DB query for all attributes
	masterAttrs, err := h.attributeRepo.FindByIDs(ctx, attrIDs)
	if err != nil {
		return nil, err
	}

	return lo.KeyBy(masterAttrs, func(attr *attributeview.AttributeView) string { return attr.ID }), nil
}

// toProductListResponse converts a list of products to HTTP responses with a single DB query for attributes
func (h *productHandler) toProductListResponse(ctx context.Context, products []*productview.ProductView) ([]httpapi.ProductResponse, error) {
	if len(products) == 0 {
		return nil, nil
	}

	// Fetch all master attributes in one query
	attrByID, err := h.fetchMasterAttributes(ctx, products)
	if err != nil {
		return nil, err
	}

	// Map all products using pre-fetched data
	return lo.Map(products, func(p *productview.ProductView, _ int) httpapi.ProductResponse {
		return *toProductResponseWithMasterData(p, attrByID)
	}), nil
}

// toProductResponse converts a single product to HTTP response (makes DB call for attributes)
func (h *productHandler) toProductResponse(ctx context.Context, p *productview.ProductView) (*httpapi.ProductResponse, error) {
	attrByID, err := h.fetchMasterAttributes(ctx, []*productview.ProductView{p})
	if err != nil {
		return nil, err
	}

	return toProductResponseWithMasterData(p, attrByID), nil
}

func (h *productHandler) GetProductById(ctx context.Context, params httpapi.GetProductByIdParams) (httpapi.GetProductByIdRes, error) { //nolint:revive // method name defined by ogen-generated interface
	q := query.GetProductByIDQuery{ID: params.ID}

	found, err := h.getByIDHandler.Handle(ctx, q)
	if errors.Is(err, mongo.ErrEntityNotFound) {
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

	// Single DB query for all attributes
	response, err := h.toProductListResponse(ctx, products)
	if err != nil {
		return nil, err
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
			//nolint:nilerr // returning HTTP 400 response instead of error
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

	// Single DB query for all attributes
	items, err := h.toProductListResponse(ctx, result.Items)
	if err != nil {
		return nil, err
	}

	return &httpapi.ProductListResponse{
		Items: items,
		Page:  result.Page,
		Size:  result.Size,
		Total: int(result.Total),
	}, nil
}

func (h *productHandler) GetProductFacets(ctx context.Context, params httpapi.GetProductFacetsParams) (httpapi.GetProductFacetsRes, error) {
	q := query.GetProductFacetsQuery{CategoryID: params.CategoryId}

	result, err := h.getFacetsHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	facets := lo.Map(result.Facets, func(f productview.AttributeFacet, _ int) httpapi.AttributeFacet {
		return httpapi.AttributeFacet{
			Slug: f.Slug,
			Values: lo.Map(f.Values, func(v productview.FacetValue, _ int) httpapi.FacetValue {
				return httpapi.FacetValue{
					Value: v.Value,
					Count: v.Count,
				}
			}),
		}
	})

	return &httpapi.FacetsResponse{
		Facets: facets,
		PriceRange: httpapi.PriceRange{
			Min: result.PriceRange.Min,
			Max: result.PriceRange.Max,
		},
	}, nil
}
