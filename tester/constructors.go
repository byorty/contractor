package tester

import (
	"github.com/byorty/contractor/common"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	NewFxTester,
	NewFxAsserterBuilder,
	NewFxPostProcessorFactory,
	fx.Annotated{
		Group:  "reporter",
		Target: NewFxStdoutReporter,
	},
	fx.Annotated{
		Group:  "worker",
		Target: NewFxWorker,
	},
	fx.Annotated{
		Group:  "worker",
		Target: NewFxWorker2,
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "eq",
				Constructor: NewEqAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "positive",
				Constructor: NewPositiveAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "min",
				Constructor: NewMinAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "max",
				Constructor: NewMaxAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "regex",
				Constructor: NewRegexAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "empty",
				Constructor: NewEmptyAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "date",
				Constructor: NewDateAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "range",
				Constructor: NewRangeAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeAsserter,
				Name:        "contains",
				Constructor: NewContainsAsserter,
			}
		},
	},
	fx.Annotated{
		Group: "post_processor_descriptor",
		Target: func(args common.Arguments) PostProcessorDescriptor {
			return PostProcessorDescriptor{
				Type: "JSON_EXTRACTOR",
				Constructor: func(config map[string]interface{}) (PostProcessor, error) {
					return NewJsonExtractorPostProcessor(args, config)
				},
			}
		},
	},
	NewFxAssertionFactory,
	NewFxEngineFactory,
	NewFxReporter2Factory,
	//fx.Annotated{
	//	Group: "assertion_descriptor",
	//	Target: func(
	//		logger common.Logger,
	//		dataCrawler common.DataCrawler,
	//		expressionFactory common.ExpressionFactory,
	//	) Asserter2Descriptor {
	//		return Asserter2Descriptor{
	//			Type: "json_contains",
	//			Constructor: func(definition interface{}) Asserter2 {
	//				return NewJsonContainsAssertion(
	//					logger,
	//					dataCrawler,
	//					expressionFactory,
	//					definition.(map[string]interface{}),
	//				)
	//			},
	//		}
	//	},
	//},
)
