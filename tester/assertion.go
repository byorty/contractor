package tester

import (
	"github.com/byorty/contractor/common"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type AssertionConstructor func(definition interface{}) Assertion2

type AssertionDescriptor struct {
	Type        string
	Constructor AssertionConstructor
}

type AssertionFactoryIn struct {
	fx.In
	Descriptors []AssertionDescriptor `group:"assertion_descriptor"`
}

type AssertionFactory interface {
	Create(name string, definition interface{}) (Assertion2, error)
}

type assertionFactory struct {
	constructors map[string]AssertionConstructor
}

func NewFxAssertionFactory(in AssertionFactoryIn) AssertionFactory {
	f := &assertionFactory{
		constructors: make(map[string]AssertionConstructor),
	}

	for _, descriptor := range in.Descriptors {
		f.constructors[descriptor.Type] = descriptor.Constructor
	}

	return f
}

func (f *assertionFactory) Create(name string, definition interface{}) (Assertion2, error) {
	constructor, ok := f.constructors[name]
	if !ok {
		return nil, errors.Errorf("assertion %s is not found", name)
	}

	return constructor(definition), nil
}

type Assertion2List struct {
	common.List[Assertion2]
}

func NewAssertion2List(assertions ...Assertion2) Assertion2List {
	return Assertion2List{common.NewListFromSlice[Assertion2](assertions...)}
}

type Assertion2 interface {
	Assert(data interface{}) (AssertionResultList, error)
}

type AssertionComparator interface {
	Compare()
}

type AssertionResultStatus int

const (
	AssertionResultStatusSuccess AssertionResultStatus = iota + 1
	AssertionResultStatusFailure
)

type AssertionResult struct {
	Status   AssertionResultStatus
	Expected string
	Actual   string
}

type AssertionResultList struct {
	common.List[AssertionResult]
}

func NewAssertionResultList(results ...AssertionResult) AssertionResultList {
	return AssertionResultList{common.NewListFromSlice[AssertionResult](results...)}
}

func (l AssertionResultList) IsPassed() bool {
	if l.Len() == 0 {
		return false
	}

	isSuccess := true
	for i := 0; i < l.Len(); i++ {
		if l.Get(i).Status == AssertionResultStatusFailure {
			isSuccess = false
			break
		}
	}

	return isSuccess
}
