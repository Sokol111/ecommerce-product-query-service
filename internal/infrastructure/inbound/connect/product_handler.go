package connect

import (
	"context"
	"errors"
	"strconv"

	"connectrpc.com/connect"
	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	productqueryv1 "github.com/Sokol111/ecommerce-product-query-service-api/gen/connect/product_query/v1"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-product-query-service/internal/application/productview"
)

type productQueryHandler struct {
	getByIDHandler   productview.GetProductByIDQueryHandler
	getRandomHandler productview.GetRandomProductsQueryHandler
	getListHandler   productview.GetListProductsQueryHandler
	getFacetsHandler productview.GetProductFacetsQueryHandler
	attributeRepo    attributeview.Repository
}

func (h *productQueryHandler) GetProductById(ctx context.Context, req *connect.Request[productqueryv1.GetProductByIdRequest]) (*connect.Response[productqueryv1.GetProductByIdResponse], error) { //nolint:revive
	found, err := h.getByIDHandler.Handle(ctx, productview.GetProductByIDQuery{ID: req.Msg.GetId()})
	if errors.Is(err, mongo.ErrEntityNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("product not found"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	proto, err := h.toProtoProduct(ctx, found)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&productqueryv1.GetProductByIdResponse{Product: proto}), nil
}

func (h *productQueryHandler) GetRandomProducts(ctx context.Context, req *connect.Request[productqueryv1.GetRandomProductsRequest]) (*connect.Response[productqueryv1.GetRandomProductsResponse], error) {
	products, err := h.getRandomHandler.Handle(ctx, productview.GetRandomProductsQuery{Count: int(req.Msg.GetCount())})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	items, err := h.toProtoProducts(ctx, products)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&productqueryv1.GetRandomProductsResponse{Items: items}), nil
}

func (h *productQueryHandler) GetProductList(ctx context.Context, req *connect.Request[productqueryv1.GetProductListRequest]) (*connect.Response[productqueryv1.GetProductListResponse], error) {
	q := productview.GetListProductsQuery{
		Page:             int(req.Msg.GetPage()),
		Size:             int(req.Msg.GetSize()),
		CategoryID:       req.Msg.CategoryId,
		Sort:             req.Msg.GetSort(),
		Order:            req.Msg.GetOrder(),
		MinPrice:         req.Msg.MinPrice,
		MaxPrice:         req.Msg.MaxPrice,
		AttributeFilters: protoToAttributeFilters(req.Msg.GetAttributeFilters()),
	}

	result, err := h.getListHandler.Handle(ctx, q)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	items, err := h.toProtoProducts(ctx, result.Items)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&productqueryv1.GetProductListResponse{
		Items: items,
		Page:  int32(result.Page), //nolint:gosec // Page originates from int32 proto field, cannot overflow
		Size:  int32(result.Size), //nolint:gosec // Size originates from int32 proto field, cannot overflow
		Total: int64(result.Total),
	}), nil
}

func (h *productQueryHandler) GetProductFacets(ctx context.Context, req *connect.Request[productqueryv1.GetProductFacetsRequest]) (*connect.Response[productqueryv1.GetProductFacetsResponse], error) {
	result, err := h.getFacetsHandler.Handle(ctx, productview.GetProductFacetsQuery{CategoryID: req.Msg.GetCategoryId()})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	facets := lo.Map(result.Facets, func(f productview.AttributeFacet, _ int) *productqueryv1.AttributeFacet {
		return &productqueryv1.AttributeFacet{
			Slug: f.Slug,
			Values: lo.Map(f.Values, func(v productview.FacetValue, _ int) *productqueryv1.FacetValue {
				return &productqueryv1.FacetValue{Value: v.Value, Count: int32(v.Count)} //nolint:gosec // Count is a small integer facet value, cannot overflow
			}),
		}
	})

	return connect.NewResponse(&productqueryv1.GetProductFacetsResponse{
		Facets: facets,
		PriceRange: &productqueryv1.PriceRange{
			Min: result.PriceRange.Min,
			Max: result.PriceRange.Max,
		},
	}), nil
}

// ==================== Helpers ====================

func (h *productQueryHandler) toProtoProduct(ctx context.Context, p *productview.ProductView) (*productqueryv1.Product, error) {
	attrByID, err := h.fetchMasterAttributes(ctx, []*productview.ProductView{p})
	if err != nil {
		return nil, err
	}
	return toProtoProductWithMasterData(p, attrByID), nil
}

func (h *productQueryHandler) toProtoProducts(ctx context.Context, products []*productview.ProductView) ([]*productqueryv1.Product, error) {
	if len(products) == 0 {
		return nil, nil
	}
	attrByID, err := h.fetchMasterAttributes(ctx, products)
	if err != nil {
		return nil, err
	}
	return lo.Map(products, func(p *productview.ProductView, _ int) *productqueryv1.Product {
		return toProtoProductWithMasterData(p, attrByID)
	}), nil
}

func (h *productQueryHandler) fetchMasterAttributes(ctx context.Context, products []*productview.ProductView) (map[string]*attributeview.AttributeView, error) {
	attrIDs := lo.Uniq(lo.FlatMap(products, func(p *productview.ProductView, _ int) []string {
		return lo.Map(p.Attributes, func(attr productview.AttributeValue, _ int) string { return attr.AttributeID })
	}))
	if len(attrIDs) == 0 {
		return make(map[string]*attributeview.AttributeView), nil
	}
	masterAttrs, err := h.attributeRepo.FindByIDs(ctx, attrIDs)
	if err != nil {
		return nil, err
	}
	return lo.KeyBy(masterAttrs, func(a *attributeview.AttributeView) string { return a.ID }), nil
}

func toProtoProductWithMasterData(p *productview.ProductView, attrByID map[string]*attributeview.AttributeView) *productqueryv1.Product {
	attrs := lo.FilterMap(p.Attributes, func(attr productview.AttributeValue, i int) (*productqueryv1.AttributeValue, bool) {
		return mapAttributeValueToProto(attr, i, attrByID)
	})
	return &productqueryv1.Product{
		Id:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		Quantity:      int32(p.Quantity), //nolint:gosec // Quantity is a product inventory count, practically bounded
		ImageId:       p.ImageID,
		SmallImageUrl: p.SmallImageURL,
		LargeImageUrl: p.LargeImageURL,
		CategoryId:    p.CategoryID,
		Attributes:    attrs,
	}
}

func mapAttributeValueToProto(attr productview.AttributeValue, index int, attrByID map[string]*attributeview.AttributeView) (*productqueryv1.AttributeValue, bool) {
	master, ok := attrByID[attr.AttributeID]
	if !ok || !master.Enabled {
		return nil, false
	}

	optionsBySlug := lo.KeyBy(master.Options, func(opt attributeview.AttributeOption) string { return opt.Slug })

	var values []*productqueryv1.AttributeValueItem

	switch string(master.Type) {
	case "single":
		if attr.OptionSlugValue != nil {
			if opt, ok := optionsBySlug[*attr.OptionSlugValue]; ok {
				values = []*productqueryv1.AttributeValueItem{{
					Slug:      attr.OptionSlugValue,
					Value:     opt.Name,
					ColorCode: opt.ColorCode,
				}}
			}
		}
	case "multiple":
		values = lo.FilterMap(attr.OptionSlugValues, func(slug string, _ int) (*productqueryv1.AttributeValueItem, bool) {
			opt, ok := optionsBySlug[slug]
			if !ok {
				return nil, false
			}
			s := slug
			return &productqueryv1.AttributeValueItem{
				Slug:      &s,
				Value:     opt.Name,
				ColorCode: opt.ColorCode,
			}, true
		})
	case "range":
		if attr.NumericValue != nil {
			values = []*productqueryv1.AttributeValueItem{{
				Value: strconv.FormatFloat(*attr.NumericValue, 'f', -1, 64),
			}}
		}
	case "boolean":
		if attr.BooleanValue != nil {
			values = []*productqueryv1.AttributeValueItem{{
				Value: strconv.FormatBool(*attr.BooleanValue),
			}}
		}
	case "text":
		if attr.TextValue != nil {
			values = []*productqueryv1.AttributeValueItem{{Value: *attr.TextValue}}
		}
	}

	return &productqueryv1.AttributeValue{
		AttributeId: attr.AttributeID,
		Slug:        attr.Slug,
		Name:        master.Name,
		Type:        stringToProtoAttributeType(string(master.Type)),
		Unit:        master.Unit,
		Role:        productqueryv1.AttributeRole_ATTRIBUTE_ROLE_SPECIFICATION,
		SortOrder:   int32(index), //nolint:gosec // index is array position in product attributes, cannot overflow
		Values:      values,
	}, true
}

func protoToAttributeFilters(filters []*productqueryv1.AttributeFilter) []productview.AttributeFilter {
	return lo.Map(filters, func(f *productqueryv1.AttributeFilter, _ int) productview.AttributeFilter {
		return productview.AttributeFilter{
			Slug:   f.GetSlug(),
			Values: f.GetValues(),
			Min:    f.Min,
			Max:    f.Max,
		}
	})
}

func stringToProtoAttributeType(s string) productqueryv1.AttributeType {
	switch s {
	case "single":
		return productqueryv1.AttributeType_ATTRIBUTE_TYPE_SINGLE
	case "multiple":
		return productqueryv1.AttributeType_ATTRIBUTE_TYPE_MULTIPLE
	case "range":
		return productqueryv1.AttributeType_ATTRIBUTE_TYPE_RANGE
	case "boolean":
		return productqueryv1.AttributeType_ATTRIBUTE_TYPE_BOOLEAN
	case "text":
		return productqueryv1.AttributeType_ATTRIBUTE_TYPE_TEXT
	default:
		return productqueryv1.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED
	}
}
