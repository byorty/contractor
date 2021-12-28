package tester

import "go.uber.org/fx"

var Constructors = fx.Provide(
	NewFxTester,
	fx.Annotated{
		Group:  "reporter",
		Target: NewFxStdoutReporter,
	},
	fx.Annotated{
		Group:  "worker",
		Target: NewFxWorker,
	},
)
