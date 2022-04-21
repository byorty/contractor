package tester

import (
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type EngineConstructor func() Engine

type EngineDescriptor struct {
	Type        string
	Constructor EngineConstructor
}

type EngineFactoryIn struct {
	fx.In
	Descriptors []EngineDescriptor `group:"engine_descriptor"`
}

type EngineFactory interface {
	Create(name string) (Engine, error)
}

type engineFactory struct {
	constructors map[string]EngineConstructor
}

func NewFxEngineFactory(in EngineFactoryIn) EngineFactory {
	f := &engineFactory{
		constructors: make(map[string]EngineConstructor),
	}

	for _, descriptor := range in.Descriptors {
		f.constructors[descriptor.Type] = descriptor.Constructor
	}

	return f
}

func (f *engineFactory) Create(name string) (Engine, error) {
	constructor, ok := f.constructors[name]
	if !ok {
		return nil, errors.Errorf("assertion %s is not found", name)
	}

	return constructor(), nil
}

type Engine interface {
	Configure(data interface{}) error
	GetTestCase2List() TestCase2List
	CreateRunner() Runner
}
