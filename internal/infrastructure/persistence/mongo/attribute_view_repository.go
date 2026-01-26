package mongo

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-product-query-service/internal/domain/attributeview"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type attributeViewRepository struct {
	*commonsmongo.GenericRepository[attributeview.AttributeView, attributeViewEntity]
	collection commonsmongo.Collection
	mapper     *attributeViewMapper
}

func newAttributeViewRepository(mongo commonsmongo.Mongo, mapper *attributeViewMapper) (attributeview.Repository, error) {
	collection := mongo.GetCollection("attribute_view")
	genericRepo, err := commonsmongo.NewGenericRepository(collection, mapper)
	if err != nil {
		return nil, err
	}

	return &attributeViewRepository{
		GenericRepository: genericRepo,
		collection:        collection,
		mapper:            mapper,
	}, nil
}

func (r *attributeViewRepository) Upsert(ctx context.Context, attribute *attributeview.AttributeView) error {
	updated, err := r.UpsertIfNewer(ctx, attribute)
	if err != nil {
		return err
	}

	if !updated {
		logger.Get(ctx).Debug("version conflict during attribute upsert",
			zap.String("id", attribute.ID),
			zap.Int("version", attribute.Version))
	}

	return nil
}

func (r *attributeViewRepository) FindByIDs(ctx context.Context, ids []string) ([]*attributeview.AttributeView, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return r.FindAllWithFilter(ctx, bson.D{{Key: "_id", Value: bson.M{"$in": ids}}}, nil)
}
