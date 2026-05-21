package mongo

import (
	"go.uber.org/fx"
)

// Module provides MongoDB infrastructure dependencies
func Module() fx.Option {
	return fx.Provide(
		newProductViewMapper,
		newProductViewRepository,
		newAttributeViewMapper,
		newAttributeViewRepository,
	)
}
