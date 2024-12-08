package routers

import "go.uber.org/fx"

var Module = fx.Provide(
	NewFarmRouter,
	MakeRouter,
)
