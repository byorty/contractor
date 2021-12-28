package converter

import "go.uber.org/fx"

var Constructors = fx.Provide(
	NewFxOa3Converter,
)
