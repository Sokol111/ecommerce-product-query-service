package mongo

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
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
	updated, err := r.UpsertIfNewer(ctx, product)
	if err != nil {
		return err
	}

	if !updated {
		logger.Get(ctx).Debug("version conflict during upsert", zap.String("id", product.ID))
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
