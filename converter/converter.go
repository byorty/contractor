package converter

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var (
	ErrUnsupportedConverterType = errors.New("unsupported spec type")
)

type Converter interface {
	GetType() common.SpecType
	Convert(ctx context.Context, arguments common.Arguments) (common.TemplateContainer, error)
}

type ContainerFxIn struct {
	fx.In
	Converters []Converter `group:"container"`
}

func NewFxContainer(in ContainerFxIn) Container {
	container := make(Container)

	for _, converter := range in.Converters {
		container[converter.GetType()] = converter
	}

	return container
}

type Container map[common.SpecType]Converter

func (c Container) Get(specType common.SpecType) (Converter, error) {
	converter, ok := c[specType]
	if !ok {
		return nil, errors.Wrap(ErrUnsupportedConverterType, string(specType))
	}

	return converter, nil
}
