package graylog

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	tc "github.com/byorty/contractor/tester/client"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	fx.Annotated{
		Group: "engine_descriptor",
		Target: func(
			loggerFactory common.LoggerFactory,
			configProviderFactory common.ConfigProviderFactory,
			graylogClient tc.GraylogClient,
		) tester.EngineDescriptor {
			return tester.EngineDescriptor{
				Type: "graylog",
				Constructor: func() tester.Engine {
					return NewFxEngine(
						loggerFactory.CreateCommonLogger(),
						configProviderFactory,
						graylogClient,
					)
				},
			}
		},
	},
)
