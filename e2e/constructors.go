package e2e

import (
	"github.com/byorty/contractor/tester/client"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	client.NewFxGraylogClient,
)
