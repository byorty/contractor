package tester

import (
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/pkg/errors"
	"github.com/spyzhov/ajson"
	"go.uber.org/fx"
)

var (
	ErrUnsupportedPostProcessor = errors.New("unsupported post processor type")
)

type PostProcessorDescriptorIn struct {
	fx.In
	Descriptors []PostProcessorDescriptor `group:"post_processor_descriptor"`
}

type PostProcessorDescriptor struct {
	Type        string
	Constructor PostProcessorConstructor
}

type PostProcessorConstructor func(config map[string]interface{}) (PostProcessor, error)

type PostProcessorFactory interface {
	Create(def common.PostProcessor) (PostProcessor, error)
}

func NewFxPostProcessorFactory(in PostProcessorDescriptorIn) PostProcessorFactory {
	factory := &postProcessorFactory{
		constructors: make(map[string]PostProcessorConstructor),
	}

	for _, descriptor := range in.Descriptors {
		factory.constructors[descriptor.Type] = descriptor.Constructor
	}

	return factory
}

type postProcessorFactory struct {
	constructors map[string]PostProcessorConstructor
}

func (f *postProcessorFactory) Create(def common.PostProcessor) (PostProcessor, error) {
	constructor, ok := f.constructors[def.Type]
	if !ok {
		return nil, errors.Wrap(ErrUnsupportedPostProcessor, def.Type)
	}

	portProcessor, err := constructor(def.Config)
	if err != nil {
		return nil, err
	}

	return portProcessor, nil
}

type PostProcessor interface {
	PostProcess(testCase *TestCase) error
}

func NewJsonExtractorPostProcessor(args common.Arguments, config map[string]interface{}) (PostProcessor, error) {
	return &jsonExtractorPostProcessor{
		arguments:    args,
		path:         fmt.Sprint(config["path"]),
		variableName: fmt.Sprint(config["variable_name"]),
	}, nil
}

type jsonExtractorPostProcessor struct {
	arguments    common.Arguments
	path         string
	variableName string
}

func (p *jsonExtractorPostProcessor) PostProcess(testCase *TestCase) error {
	root, err := ajson.Unmarshal(testCase.ActualResult.Buf)
	if err != nil {
		return err
	}

	nodes, err := root.JSONPath(p.path)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		val, err := node.Value()
		if err != nil {
			return err
		}

		p.arguments.Variables[p.variableName] = fmt.Sprint(val)
	}

	return nil
}
