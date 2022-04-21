package tester

import (
	"github.com/byorty/contractor/common"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type Asserter2Constructor func(definition interface{}) Asserter2

type Asserter2Descriptor struct {
	Type        string
	Constructor Asserter2Constructor
}

type Assertion2FactoryIn struct {
	fx.In
	Descriptors []Asserter2Descriptor `group:"assertion_descriptor"`
}

type Assertion2Factory interface {
	Create(name string, definition interface{}) (Asserter2, error)
}

type assertionFactory struct {
	constructors map[string]Asserter2Constructor
}

func NewFxAssertionFactory(in Assertion2FactoryIn) Assertion2Factory {
	f := &assertionFactory{
		constructors: make(map[string]Asserter2Constructor),
	}

	for _, descriptor := range in.Descriptors {
		f.constructors[descriptor.Type] = descriptor.Constructor
	}

	return f
}

func (f *assertionFactory) Create(name string, definition interface{}) (Asserter2, error) {
	constructor, ok := f.constructors[name]
	if !ok {
		return nil, errors.Errorf("assertion %s is not found", name)
	}

	return constructor(definition), nil
}

type Asserter2List struct {
	common.List[Asserter2]
}

func NewAsserter2List(assertions ...Asserter2) Asserter2List {
	return Asserter2List{common.NewListFromSlice[Asserter2](assertions...)}
}

type Asserter2 interface {
	Assert(data interface{}) AssertionResultList
}

type AssertionResultStatus int

const (
	AssertionResultStatusSuccess AssertionResultStatus = iota + 1
	AssertionResultStatusFailure
)

type AssertionResult struct {
	Name     string
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
