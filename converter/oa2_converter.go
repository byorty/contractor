package converter

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func NewFxOa2Converter() Converter {
	return &oa2Converter{
		oaConverter: &oaConverter{
			container: make(common.TemplateContainer),
		},
	}
}

type oa2Converter struct {
	*oaConverter
}

func (c *oa2Converter) GetType() common.SpecType {
	return common.SpecTypeOA2
}

func (c *oa2Converter) Convert(ctx context.Context, arguments common.Arguments) (common.TemplateContainer, error) {
	var specV2 openapi2.T
	err := c.readAndUnmarshal(arguments.SpecFilename, &specV2)
	if err != nil {
		return nil, err
	}

	specV3, err := openapi2conv.ToV3(&specV2)
	if err != nil {
		return nil, err
	}

	return c.convert(ctx, arguments, specV3)
}
