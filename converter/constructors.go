package converter

import "go.uber.org/fx"

var Constructors = fx.Provide(
	NewFxContainer,
	fx.Annotated{
		Group:  "container",
		Target: NewFxOa2Converter,
	},
	fx.Annotated{
		Group:  "container",
		Target: NewFxOa3Converter,
	},
)
