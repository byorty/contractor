package common

import (
	"fmt"
	"github.com/antonmedv/expr"
	"go.uber.org/fx"
)

type ExpressionType string

const (
	ExpressionTypeAsserter  ExpressionType = "asserter"
	ExpressionTypeGenerator ExpressionType = "generator"
)

const (
	ExpressionDateLayoutRFC3339     string = "RFC3339"
	ExpressionDateLayoutRFC3339Nano string = "RFC3339NANO"
)

type ExpressionFactoryIn struct {
	fx.In
	Descriptors []ExpressionDescriptor `group:"expression_descriptor"`
}

type ExpressionDescriptor struct {
	Type        ExpressionType
	Name        string
	Constructor interface{}
}

type ExpressionFactory interface {
	Create(exprType ExpressionType, exprName interface{}) (interface{}, error)
}

func NewFxExpressionFactory(in ExpressionFactoryIn) ExpressionFactory {
	f := &expressionFactory{
		constructors: make(map[string]interface{}),
	}

	for _, descriptor := range in.Descriptors {
		f.constructors[fmt.Sprintf("%s_%v", descriptor.Type, descriptor.Name)] = descriptor.Constructor
	}

	return f
}

type expressionFactory struct {
	constructors map[string]interface{}
}

func (f *expressionFactory) Create(exprType ExpressionType, exprName interface{}) (interface{}, error) {
	output, err := expr.Eval(fmt.Sprintf("%s_%v", exprType, exprName), f.constructors)
	if err != nil {
		return nil, err
	}

	return output, nil
}
