package reporter

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	fx.Annotated{
		Group: "reporter_descriptor",
		Target: func(
			loggerFactory common.LoggerFactory,
		) tester.Reporter2Descriptor {
			return tester.Reporter2Descriptor{
				Type: "console",
				Constructor: func(definition interface{}) tester.Reporter2 {
					return NewConsole(loggerFactory)
				},
			}
		},
	},
)
