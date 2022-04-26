package graylog

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	fx.Annotated{
		Group: "engine_descriptor",
		Target: func(
			ctx context.Context,
			loggerFactory common.LoggerFactory,
			configProviderFactory common.ConfigProviderFactory,
		) tester.EngineDescriptor {
			return tester.EngineDescriptor{
				Type: "graylog",
				Constructor: func() tester.Engine {
					return NewFxEngine(
						ctx,
						loggerFactory.CreateCommonLogger(),
						configProviderFactory,
					)
				},
			}
		},
	},
)
