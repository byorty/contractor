package common

import "github.com/pkg/errors"

type Descriptor[T any] struct {
	Type        string
	Constructor T
}

type Factory[T any] interface {
	Create(t string) (T, error)
}

func NewFactory[T any](descriptors []Descriptor[T]) Factory[T] {
	f := &factory[T]{
		constructors: NewMap[string, T](),
	}

	for _, descriptor := range descriptors {
		f.constructors.Set(descriptor.Type, descriptor.Constructor)
	}

	return f
}

type factory[T any] struct {
	constructors Map[string, T]
}

func (f factory[T]) Create(t string) (T, error) {
	constructor, ok := f.constructors.Get(t)
	if !ok {
		return nil, errors.Errorf("type %s is not found", t)
	}

	return constructor, nil
}
