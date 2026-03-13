package mongo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type productViewRepository struct {
	*commonsmongo.GenericRepository[productview.ProductView, productViewEntity]
	collection *mongodriver.Collection
	mapper     *productViewMapper
}

func newProductViewRepository(mongo commonsmongo.Mongo, mapper *productViewMapper) (productview.Repository, error) {
	collection := mongo.GetCollection("product_view")
	genericRepo, err := commonsmongo.NewGenericRepository(collection, mapper)
	if err != nil {
		return nil, err
	}

	return &productViewRepository{
		GenericRepository: genericRepo,
		collection:        collection,
		mapper:            mapper,
	}, nil
}

func (r *productViewRepository) Upsert(ctx context.Context, product *productview.ProductView) error {
	entity := r.mapper.ToEntity(product)

	filter := bson.M{
		"_id":     entity.ID,
		"version": bson.M{"$lt": entity.Version},
	}

	// Replace document, but preserve image URLs if imageId hasn't changed.
	// This prevents "flickering" when product is updated but image stays the same,
	// since image-service only publishes ProductImagePromoted on promotion, not on every product update.
	update := bson.A{
		bson.M{
			"$replaceWith": bson.M{
				"$mergeObjects": bson.A{
					entity,
					bson.M{
						"smallImageUrl": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$imageId", entity.ImageID}},
								"then": "$smallImageUrl",
								"else": nil,
							},
						},
						"largeImageUrl": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$imageId", entity.ImageID}},
								"then": "$largeImageUrl",
								"else": nil,
							},
						},
					},
				},
			},
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert product view: %w", err)
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		logger.Get(ctx).Debug("version conflict during upsert", zap.String("id", product.ID))
	}

	return nil
}

func (r *productViewRepository) UpdateImageURLs(ctx context.Context, productID, imageID, smallImageURL, largeImageURL string) error {
	// Only update if the document's imageId matches - prevents overwriting with stale data
	// when product was already updated with a different image
	filter := bson.D{
		{Key: "_id", Value: productID},
		{Key: "imageId", Value: imageID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "smallImageUrl", Value: smallImageURL},
			{Key: "largeImageUrl", Value: largeImageURL},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update image URLs: %w", err)
	}

	if result.MatchedCount == 0 {
		logger.Get(ctx).Debug("skipped image URLs update - imageId mismatch or product not found",
			zap.String("productId", productID),
			zap.String("imageId", imageID))
	}

	return nil
}

func (r *productViewRepository) FindRandom(ctx context.Context, count int) ([]*productview.ProductView, error) {
	pipeline := mongodriver.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "enabled", Value: true}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: count}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate random products: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }() //nolint:errcheck // cursor.Close error is not actionable

	var entities []productViewEntity
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return lo.Map(entities, func(e productViewEntity, _ int) *productview.ProductView {
		return r.mapper.ToDomain(&e)
	}), nil
}

func (r *productViewRepository) FindList(ctx context.Context, query productview.ListQuery) (*commonsmongo.PageResult[productview.ProductView], error) {
	filter := r.buildListFilter(query)

	var sortBson bson.D
	if query.Sort != "" {
		sortOrder := 1 // asc
		if query.Order == "desc" {
			sortOrder = -1
		}
		sortBson = bson.D{{Key: query.Sort, Value: sortOrder}}
	}

	opts := commonsmongo.QueryOptions{
		Filter: filter,
		Page:   query.Page,
		Size:   query.Size,
		Sort:   sortBson,
	}

	return r.FindWithOptions(ctx, opts)
}

func (r *productViewRepository) buildListFilter(query productview.ListQuery) bson.D {
	filter := bson.D{{Key: "enabled", Value: true}}

	if query.CategoryID != nil {
		filter = append(filter, bson.E{Key: "categoryId", Value: *query.CategoryID})
	}

	filter = r.appendPriceFilter(filter, query.MinPrice, query.MaxPrice)
	filter = r.appendAttributeFilters(filter, query.AttributeFilters)

	return filter
}

func (r *productViewRepository) appendPriceFilter(filter bson.D, minPrice, maxPrice *float64) bson.D {
	if minPrice == nil && maxPrice == nil {
		return filter
	}

	priceFilter := bson.M{}
	if minPrice != nil {
		priceFilter["$gte"] = *minPrice
	}
	if maxPrice != nil {
		priceFilter["$lte"] = *maxPrice
	}
	return append(filter, bson.E{Key: "price", Value: priceFilter})
}

func coerceFilterValue(s string) any {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s
}

func (r *productViewRepository) appendAttributeFilters(filter bson.D, attrFilters []productview.AttributeFilter) bson.D {
	for _, attrFilter := range attrFilters {
		attrKey := "attrs." + attrFilter.Slug

		if len(attrFilter.Values) > 0 {
			// Coerce string values to native types (bool/numeric) for proper MongoDB type matching
			values := lo.Map(attrFilter.Values, func(v string, _ int) any {
				return coerceFilterValue(v)
			})
			filter = append(filter, bson.E{Key: attrKey, Value: bson.M{"$in": values}})
		} else if attrFilter.Min != nil || attrFilter.Max != nil {
			// For range type: numeric comparison
			rangeFilter := bson.M{}
			if attrFilter.Min != nil {
				rangeFilter["$gte"] = *attrFilter.Min
			}
			if attrFilter.Max != nil {
				rangeFilter["$lte"] = *attrFilter.Max
			}
			filter = append(filter, bson.E{Key: attrKey, Value: rangeFilter})
		}
	}
	return filter
}

func (r *productViewRepository) FindFacets(ctx context.Context, categoryID string) (*productview.FacetsResult, error) {
	// Match only enabled products in the given category that have attrs
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "enabled", Value: true},
			{Key: "categoryId", Value: categoryID},
			{Key: "attrs", Value: bson.M{"$exists": true, "$ne": bson.M{}}},
		}},
	}

	// Use $facet to compute attribute facets and price range in a single aggregation
	facetStage := bson.D{
		{Key: "$facet", Value: bson.D{
			{Key: "attrFacets", Value: bson.A{
				// Convert attrs map to array of {k, v} pairs
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "attrsArray", Value: bson.M{"$objectToArray": "$attrs"}},
				}}},
				// Unwind to get one doc per attribute
				bson.D{{Key: "$unwind", Value: "$attrsArray"}},
				// Unwind arrays (for multiple type) — leaves scalars intact
				bson.D{{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$attrsArray.v"},
					{Key: "preserveNullAndEmptyArrays", Value: false},
				}}},
				// Group by attribute slug + value, count products
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "slug", Value: "$attrsArray.k"},
						{Key: "value", Value: "$attrsArray.v"},
					}},
					{Key: "count", Value: bson.M{"$sum": 1}},
				}}},
				// Sort by slug then count (descending)
				bson.D{{Key: "$sort", Value: bson.D{
					{Key: "_id.slug", Value: 1},
					{Key: "count", Value: -1},
				}}},
				// Group by attribute slug to collect all values
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$_id.slug"},
					{Key: "values", Value: bson.M{"$push": bson.D{
						{Key: "value", Value: "$_id.value"},
						{Key: "count", Value: "$count"},
					}}},
				}}},
				// Sort by slug for consistent ordering
				bson.D{{Key: "$sort", Value: bson.D{
					{Key: "_id", Value: 1},
				}}},
			}},
			{Key: "priceRange", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "min", Value: bson.M{"$min": "$price"}},
					{Key: "max", Value: bson.M{"$max": "$price"}},
				}}},
			}},
		}},
	}

	pipeline := mongodriver.Pipeline{matchStage, facetStage}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate facets: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }() //nolint:errcheck // cursor.Close error is not actionable

	var results []facetsAggregationResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode facets: %w", err)
	}

	if len(results) == 0 {
		return &productview.FacetsResult{}, nil
	}

	return mapFacetsResult(&results[0]), nil
}

// facetsAggregationResult represents the MongoDB $facet output structure
type facetsAggregationResult struct {
	AttrFacets []attrFacetResult `bson:"attrFacets"`
	PriceRange []priceResult     `bson:"priceRange"`
}

type attrFacetResult struct {
	Slug   string             `bson:"_id"`
	Values []facetValueResult `bson:"values"`
}

type facetValueResult struct {
	Value any `bson:"value"`
	Count int `bson:"count"`
}

type priceResult struct {
	Min float64 `bson:"min"`
	Max float64 `bson:"max"`
}

func mapFacetsResult(result *facetsAggregationResult) *productview.FacetsResult {
	facets := make([]productview.AttributeFacet, len(result.AttrFacets))
	for i, af := range result.AttrFacets {
		values := make([]productview.FacetValue, len(af.Values))
		for j, v := range af.Values {
			values[j] = productview.FacetValue{
				Value: fmt.Sprintf("%v", v.Value),
				Count: v.Count,
			}
			// For string values (single/multiple type), use value as slug
			if strVal, ok := v.Value.(string); ok {
				values[j].Slug = strVal
			}
		}
		facets[i] = productview.AttributeFacet{
			Slug:   af.Slug,
			Values: values,
		}
	}

	var priceRange productview.PriceRange
	if len(result.PriceRange) > 0 {
		priceRange = productview.PriceRange{
			Min: result.PriceRange[0].Min,
			Max: result.PriceRange[0].Max,
		}
	}

	return &productview.FacetsResult{
		Facets:     facets,
		PriceRange: priceRange,
	}
}
