package tester

import (
	"github.com/byorty/contractor/common"
	"go.uber.org/fx"
	"time"
)

var Constructors = fx.Provide(
	NewFxTester,
	NewFxAsserterBuilder,
	fx.Annotated{
		Group:  "reporter",
		Target: NewFxStdoutReporter,
	},
	fx.Annotated{
		Group:  "worker",
		Target: NewFxWorker,
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeAsserter,
				Name: "eq",
				Constructor: func(expected interface{}) Asserter {
					return &eqAsserter{
						expected: expected,
					}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeAsserter,
				Name: "positive",
				Constructor: func() Asserter {
					return &minAsserter{
						expected: 1,
					}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeAsserter,
				Name: "min",
				Constructor: func(expected float64) Asserter {
					return &minAsserter{
						expected: expected,
					}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeAsserter,
				Name: "regex",
				Constructor: func(expr string) Asserter {
					return &regexAsserter{
						expected: expr,
					}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeAsserter,
				Name: "empty",
				Constructor: func() Asserter {
					return &emptyAsserter{}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeAsserter,
				Name: "date",
				Constructor: func(layout string) Asserter {
					var expected string
					switch layout {
					case "RFC3339":
						expected = time.RFC3339
					case "RFC3339NANO":
						expected = time.RFC3339Nano
					default:
						expected = layout
					}

					return &dateAsserter{
						expected: expected,
					}
				},
			}
		},
	},
)
