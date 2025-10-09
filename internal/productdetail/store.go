package productdetail

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var errEntityNotFound = errors.New("entity not found in database")

type Store interface {
	Upsert(ctx context.Context, id string, name string, description string, price float32, quantity int, version int, enabled bool) error

	GetById(ctx context.Context, id string) (*model.ProductDTO, error)

	GetRandomProducts(ctx context.Context, amount int) ([]*model.ProductDTO, error)
}

type store struct {
	wrapper *mongo.CollectionWrapper[collection]
}

func newStore(wrapper *mongo.CollectionWrapper[collection]) Store {
	return &store{wrapper}
}

func (s *store) Upsert(ctx context.Context, id string, name string, description string, price float32, quantity int, version int, enabled bool) error {
	filter := bson.M{
		"_id":     id,
		"version": bson.M{"$lt": version},
	}
	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"description": description,
			"enabled":     enabled,
			"version":     version,
			"price":       price,
			"quantity":    quantity,
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

func (s *store) GetById(ctx context.Context, id string) (*model.ProductDTO, error) {
	result := s.wrapper.Coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	var e entity
	err := result.Decode(&e)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, fmt.Errorf("failed to get product [%v]: %w", id, errEntityNotFound)
		}
		return nil, fmt.Errorf("failed to get product [%v]: decode error: %w", id, err)
	}

	return toDTO(&e), nil
}

func (s *store) GetRandomProducts(ctx context.Context, amount int) ([]*model.ProductDTO, error) {
	pipeline := mongodriver.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "enabled", Value: true}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: amount}}}},
	}
	cursor, err := s.wrapper.Coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate random products, amount [%v]: %w", amount, err)
	}

	defer cursor.Close(ctx)

	var entities []entity
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	dtos := make([]*model.ProductDTO, 0, len(entities))

	for i := range entities {
		dtos = append(dtos, toDTO(&entities[i]))
	}

	return dtos, nil
}

func toDTO(e *entity) *model.ProductDTO {
	return &model.ProductDTO{ID: e.ID, Name: e.Name, Price: e.Price, Quantity: e.Quantity}
}

func (s *store) log(ctx context.Context) *zap.Logger {
	return logger.FromContext(ctx).With(zap.String("component", "product-detail-store"))
}
