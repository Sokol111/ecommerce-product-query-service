package categorylist

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-category-query-service/internal/model"
	"github.com/Sokol111/ecommerce-commons/pkg/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Store interface {
	Upsert(ctx context.Context, id string, name string, price float32, quantity int, version int, enabled bool) error

	GetAllEnabled(ctx context.Context) (*model.CategoryListViewDTO, error)
}

type store struct {
	wrapper *mongo.CollectionWrapper[collection]
	logger  *zap.Logger
}

func newStore(wrapper *mongo.CollectionWrapper[collection], logger *zap.Logger) Store {
	return &store{wrapper, logger.With(zap.String("component", "product-detail-store"))}
}

func (s *store) Upsert(ctx context.Context, id string, name string, price float32, quantity int, version int, enabled bool) error {
	filter := bson.M{
		"_id":     id,
		"version": bson.M{"$lt": version},
	}
	update := bson.M{
		"$set": bson.M{
			"name":     name,
			"enabled":  enabled,
			"version":  version,
			"price":    price,
			"quantity": quantity,
		},
	}
	opts := options.Update().SetUpsert(true)
	result, err := s.wrapper.Coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert category: %w", err)
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		s.log(ctx).Debug("version conflict", zap.String("id", id))
	}
	return nil
}

func (s *store) GetAllEnabled(ctx context.Context) (*model.CategoryListViewDTO, error) {
	cursor, err := s.wrapper.Coll.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return nil, fmt.Errorf("failed to get active categories: %w", err)
	}
	defer cursor.Close(ctx)

	var categories []model.CategoryDTO
	for cursor.Next(ctx) {
		var doc struct {
			ID       string  `bson:"_id"`
			Name     string  `bson:"name"`
			Price    float32 `bson:"price"`
			Quantity int     `bson:"quantity"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode category: %w", err)
		}
		categories = append(categories, model.CategoryDTO{ID: doc.ID, Name: doc.Name, Price: doc.Price, Quantity: doc.Quantity})
	}
	return &model.ProductDTO{Categories: categories}, nil
}

func (s *store) log(ctx context.Context) *zap.Logger {
	return logger.CombineLogger(s.logger, ctx)
}
