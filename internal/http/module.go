package http //nolint:revive // intentional package name

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			newProductHandler,
		),
	)
}
