package mocker

import (
	"github.com/byorty/contractor/common"
	"go.uber.org/fx"
	"time"
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
				Type: common.ExpressionTypeGenerator,
				Name: "eq",
				Constructor: func(value interface{}) Generator {
					return &eqGenerator{value: value}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeGenerator,
				Name: "positive",
				Constructor: func() Generator {
					return &minGenerator{min: 1}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeGenerator,
				Name: "min",
				Constructor: func(min int) Generator {
					return &minGenerator{min: min}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeGenerator,
				Name: "regex",
				Constructor: func(expr string) Generator {
					return &regexGenerator{expr: expr}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeGenerator,
				Name: "empty",
				Constructor: func() Generator {
					return &emptyGenerator{}
				},
			}
		},
	},
	fx.Annotated{
		Group: "expression_descriptor",
		Target: func() common.ExpressionDescriptor {
			return common.ExpressionDescriptor{
				Type: common.ExpressionTypeGenerator,
				Name: "date",
				Constructor: func(layout string) Generator {
					switch layout {
					case "RFC3339":
						layout = time.RFC3339
					case "RFC3339NANO":
						layout = time.RFC3339Nano
					}

					return &dateGenerator{layout: layout}
				},
			}
		},
	},
)
