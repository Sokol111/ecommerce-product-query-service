package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/productview"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type productViewRepository struct {
	collection commonsmongo.Collection
	mapper     *productViewMapper
}

func newProductViewRepository(mongo commonsmongo.Mongo, mapper *productViewMapper) productview.Repository {
	return &productViewRepository{
		collection: mongo.GetCollection("product_detail"),
		mapper:     mapper,
	}
}

func (r *productViewRepository) Upsert(ctx context.Context, product *productview.ProductView) error {
	entity := r.mapper.ToEntity(product)

	filter := bson.M{
		"_id":     entity.ID,
		"version": bson.M{"$lt": entity.Version},
	}

	update := bson.M{
		"$set": bson.M{
			"name":        entity.Name,
			"description": entity.Description,
			"enabled":     entity.Enabled,
			"version":     entity.Version,
			"price":       entity.Price,
			"quantity":    entity.Quantity,
			"imageId":     entity.ImageID,
			"createdAt":   entity.CreatedAt,
			"modifiedAt":  entity.ModifiedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert product view: %w", err)
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		logger.FromContext(ctx).Debug("version conflict during upsert", zap.String("id", product.ID))
	}

	return nil
}

func (r *productViewRepository) FindByID(ctx context.Context, id string) (*productview.ProductView, error) {
	var entity productViewEntity

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, persistence.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to find product view by id: %w", err)
	}

	return r.mapper.ToDomain(&entity), nil
}

func (r *productViewRepository) FindRandom(ctx context.Context, amount int) ([]*productview.ProductView, error) {
	pipeline := mongodriver.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "enabled", Value: true}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: amount}}}},
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
