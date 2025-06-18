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
	Upsert(ctx context.Context, id string, name string, price float32, quantity int, version int, enabled bool) error

	GetById(ctx context.Context, id string) (*model.ProductDTO, error)
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

func (s *store) GetById(ctx context.Context, id string) (*model.ProductDTO, error) {
	result := s.wrapper.Coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	var doc struct {
		ID       string  `bson:"_id"`
		Name     string  `bson:"name"`
		Price    float32 `bson:"price"`
		Quantity int     `bson:"quantity"`
	}
	err := result.Decode(&doc)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, fmt.Errorf("failed to get product [%v]: %w", id, errEntityNotFound)
		}
		return nil, fmt.Errorf("failed to get product [%v]: decode error: %w", id, err)
	}

	return &model.ProductDTO{
		ID:       doc.ID,
		Name:     doc.Name,
		Price:    doc.Price,
		Quantity: doc.Quantity,
	}, nil
}

func (s *store) log(ctx context.Context) *zap.Logger {
	return logger.CombineLogger(s.logger, ctx)
}
