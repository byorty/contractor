package mocker

import (
	"github.com/byorty/contractor/common"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	NewFxMocker,
	fx.Annotated{
		Group:  "worker",
		Target: NewFxWorker,
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "eq",
				Constructor: NewEqGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "positive",
				Constructor: NewPositiveGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "min",
				Constructor: NewMinGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "max",
				Constructor: NewMaxGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "range",
				Constructor: NewRangeGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "regex",
				Constructor: NewRegexGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "empty",
				Constructor: NewEmptyGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "date",
				Constructor: NewDateGenerator,
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type:        common.ExpressionTypeGenerator,
				Name:        "contains",
				Constructor: NewContainsGenerator,
			}
		},
	},
)
