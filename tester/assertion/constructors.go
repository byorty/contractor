package assertion

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	fx.Annotated{
		Group: "assertion_descriptor",
		Target: func(
			loggerFactory common.LoggerFactory,
			dataCrawler common.DataCrawler,
			expressionFactory common.ExpressionFactory,
		) tester.Asserter2Descriptor {
			return tester.Asserter2Descriptor{
				Type: "json_contains",
				Constructor: func(definition interface{}) tester.Asserter2 {
					return NewJsonContains(
						loggerFactory.CreateCommonLogger(),
						dataCrawler,
						expressionFactory,
						definition,
					)
				},
			}
		},
	},
)
