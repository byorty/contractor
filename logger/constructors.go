package logger

import "go.uber.org/fx"

var Constructors = fx.Provide(
	NewFxLogger,
)
