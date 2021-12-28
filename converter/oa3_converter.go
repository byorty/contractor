package converter

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/getkin/kin-openapi/openapi3"
	"strconv"
)

func NewFxOa3Converter(
	mediaConverter common.MediaConverter,
) Converter {
	return &oa3Converter{
		mediaConverter: mediaConverter,
	}
}

type oa3Converter struct {
	mediaConverter common.MediaConverter
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

	container := make(common.TemplateContainer)
	for pathName, pathRef := range spec.Paths {
		for httpMethod, operation := range pathRef.Operations() {
			for _, parameterRef := range operation.Parameters {
				parameter := parameterRef.Value

				for exampleName, exampleRef := range parameter.Examples {
					template, ok := container[exampleName]
					if !ok {
						template = common.Template{
							BaseUrl:           arguments.BaseUrl,
							UID:               operation.OperationID,
							Path:              pathName,
							Method:            httpMethod,
							PathParams:        make(map[string]interface{}),
							QueryParams:       make(map[string]interface{}),
							HeaderParams:      make(map[string]interface{}),
							CookieParams:      make(map[string]interface{}),
							ExpectedResponses: make(map[int]map[string]interface{}),
						}

						container[exampleName] = template
					}

					switch parameter.In {
					case openapi3.ParameterInPath:
						container[exampleName].PathParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInQuery:
						container[exampleName].QueryParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInHeader:
						container[exampleName].HeaderParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInCookie:
						container[exampleName].CookieParams[parameter.Name] = exampleRef.Value.Value
					}
				}
			}

			for statusCodeName, responseRef := range operation.Responses {
				statusCode, err := strconv.Atoi(statusCodeName)
				if err != nil {
					continue
				}

				for mediaTypeName, mediaType := range responseRef.Value.Content {
					for exampleName, exampleRef := range mediaType.Examples {
						template := container[exampleName]

						if template.ExpectedResponses == nil {
							continue
						}

						_, ok := template.ExpectedResponses[statusCode]
						if !ok {
							template.ExpectedResponses[statusCode] = make(map[string]interface{})
						}

						template.ExpectedResponses[statusCode][mediaTypeName] = exampleRef.Value.Value
					}
				}
			}
		}
	}

	return container, nil
}
