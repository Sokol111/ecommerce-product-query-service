package productdetail

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"
)

func NewCategoryListViewModule() fx.Option {
	return fx.Provide(
		provideCollection,
		newStore,
		newService,
	)
}

func provideCollection(lc fx.Lifecycle, m mongo.Mongo) (*mongo.CollectionWrapper[collection], error) {
	wrapper := &mongo.CollectionWrapper[collection]{}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wrapper.Coll = m.GetCollection("product_detail")
			err := m.CreateSimpleIndex(ctx, "product_detail", bson.D{{Key: "enabled", Value: 1}})
			if err != nil {
				return err
			}
			return nil
		},
	})
	return wrapper, nil
}
