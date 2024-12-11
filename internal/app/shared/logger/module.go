package shared

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewLogger,
)
