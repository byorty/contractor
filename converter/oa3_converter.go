package converter

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/getkin/kin-openapi/openapi3"
)

func NewFxOa3Converter() Converter {
	return &oa3Converter{
		oaConverter: &oaConverter{
			container: make(common.TemplateContainer),
		},
	}
}

type oa3Converter struct {
	*oaConverter
}

func (c *oa3Converter) GetType() common.SpecType {
	return common.SpecTypeOA3
}

func (c *oa3Converter) Convert(ctx context.Context, arguments common.Arguments) (common.TemplateContainer, error) {
	loader := &openapi3.Loader{Context: ctx}
	spec, err := loader.LoadFromFile(arguments.SpecFilename)
	if err != nil {
		return nil, err
	}

	err = spec.Validate(ctx)
	if err != nil {
		return nil, err
	}

	return c.convert(ctx, arguments, spec)
}
