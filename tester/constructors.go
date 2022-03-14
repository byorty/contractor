package tester

import (
	"github.com/byorty/contractor/common"
	"go.uber.org/fx"
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
)
