package tester

import (
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type Reporter2Constructor func(definition interface{}) Reporter2

type Reporter2Descriptor struct {
	Type        string
	Constructor Reporter2Constructor
}

type Reporter2FactoryIn struct {
	fx.In
	Descriptors []Reporter2Descriptor `group:"reporter_descriptor"`
}

type Reporter2Factory interface {
	Create(name string, definition interface{}) (Reporter2, error)
}

type reporterFactory struct {
	constructors map[string]Reporter2Constructor
}

func NewFxReporter2Factory(in Reporter2FactoryIn) Reporter2Factory {
	f := &reporterFactory{
		constructors: make(map[string]Reporter2Constructor),
	}

	for _, descriptor := range in.Descriptors {
		f.constructors[descriptor.Type] = descriptor.Constructor
	}

	return f
}

func (f *reporterFactory) Create(name string, definition interface{}) (Reporter2, error) {
	constructor, ok := f.constructors[name]
	if !ok {
		return nil, errors.Errorf("assertion %s is not found", name)
	}

	return constructor(definition), nil
}

type Reporter2 interface {
	Report(report RunnerReport)
}
