package mocker

import "go.uber.org/fx"

var Constructors = fx.Provide(
	NewFxMocker,
	fx.Annotated{
		Group:  "worker",
		Target: NewFxWorker,
	},
)
