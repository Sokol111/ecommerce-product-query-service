package mongo

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type productViewRepository struct {
	*commonsmongo.GenericRepository[productview.ProductView, productViewEntity]
	collection commonsmongo.Collection
	mapper     *productViewMapper
}

func newProductViewRepository(mongo commonsmongo.Mongo, mapper *productViewMapper) (productview.Repository, error) {
	collection := mongo.GetCollection("product_detail")
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

	// Replace document, but preserve imageUrl if imageId hasn't changed.
	// This prevents "flickering" when product is updated but image stays the same,
	// since image-service only publishes ProductImagePromoted on promotion, not on every product update.
	update := bson.A{
		bson.M{
			"$replaceWith": bson.M{
				"$mergeObjects": bson.A{
					entity,
					bson.M{
						"imageUrl": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$imageId", entity.ImageID}},
								"then": "$imageUrl",
								"else": nil,
							},
						},
					},
				},
			},
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert product view: %w", err)
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		logger.Get(ctx).Debug("version conflict during upsert", zap.String("id", product.ID))
	}

	return nil
}

func (r *productViewRepository) UpdateImageURL(ctx context.Context, productID, imageID, imageURL string) error {
	// Only update if the document's imageId matches - prevents overwriting with stale data
	// when product was already updated with a different image
	filter := bson.D{
		{Key: "_id", Value: productID},
		{Key: "imageId", Value: imageID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "imageUrl", Value: imageURL},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update image URL: %w", err)
	}

	if result.MatchedCount == 0 {
		logger.Get(ctx).Debug("skipped image URL update - imageId mismatch or product not found",
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
	defer cursor.Close(ctx)

	var entities []productViewEntity
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	views := make([]*productview.ProductView, 0, len(entities))
	for i := range entities {
		views = append(views, r.mapper.ToDomain(&entities[i]))
	}

	return views, nil
}

func (r *productViewRepository) FindList(ctx context.Context, query productview.ListQuery) (*commonsmongo.PageResult[productview.ProductView], error) {
	filter := bson.D{{Key: "enabled", Value: true}}

	if query.CategoryID != nil {
		filter = append(filter, bson.E{Key: "categoryId", Value: *query.CategoryID})
	}

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
