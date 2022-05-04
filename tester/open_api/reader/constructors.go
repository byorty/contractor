package reader

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	NewFxFactory,
	fx.Annotated{
		Group: "open_api_reader",
		Target: common.Descriptor[Constructor]{
			Type: "oa2",
			Constructor: func() Reader {
				return &reader{
					ctx: nil,
					readerFunc: func(ctx context.Context, filename string) (*openapi3.T, error) {
						//var specV2 openapi2.T
						//err := c.readAndUnmarshal(filename, &specV2)
						//if err != nil {
						//	return nil, err
						//}
						//
						//specV3, err := openapi2conv.ToV3(&specV2)
						//if err != nil {
						//	return nil, err
						//}
					},
				}
			},
		},
	},
	fx.Annotated{
		Group: "open_api_reader",
		Target: common.Descriptor[Constructor]{
			Type: "oa3",
			Constructor: func() Reader {
				return &reader{
					ctx:        nil,
					readerFunc: nil,
				}
			},
		},
	},
)
