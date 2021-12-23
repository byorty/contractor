package converter

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/logger"
	"github.com/getkin/kin-openapi/openapi3"
	"strconv"
)

func NewFxOa3Converter(
	logger *logger.Logger,
) Converter {
	return &oa3Converter{
		logger: logger,
	}
}

type oa3Converter struct {
	logger *logger.Logger
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
			c.logger.PrintH1("Operation: %s", operation.OperationID)
			c.logger.PrintH2("Path: %s", pathName)
			collection, ok := container[operation.OperationID]
			if !ok {
				collection = make(map[string]common.Template)
				container[operation.OperationID] = collection
			}

			for _, parameterRef := range operation.Parameters {
				parameter := parameterRef.Value

				for exampleName, exampleRef := range parameter.Examples {
					template, ok := collection[exampleName]
					if !ok {
						template = common.Template{
							BaseUrl:           arguments.BaseUrl,
							Path:              pathName,
							Method:            httpMethod,
							PathParams:        make(map[string]interface{}),
							QueryParams:       make(map[string]interface{}),
							HeaderParams:      make(map[string]interface{}),
							CookieParams:      make(map[string]interface{}),
							ExpectedResponses: make(map[int]map[string]interface{}),
						}

						collection[exampleName] = template
					}

					switch parameter.In {
					case openapi3.ParameterInPath:
						collection[exampleName].PathParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInQuery:
						collection[exampleName].QueryParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInHeader:
						collection[exampleName].HeaderParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInCookie:
						collection[exampleName].CookieParams[parameter.Name] = exampleRef.Value.Value
					}

				}
			}

			for example, template := range collection {
				c.logger.PrintH2("Example: %s", example)
				c.logger.PrintH2Parameters("Header Parameters", template.HeaderParams)
				c.logger.PrintH2Parameters("Path Parameters", template.PathParams)
				c.logger.PrintH2Parameters("Query Parameters", template.QueryParams)
				c.logger.PrintH2Parameters("Cookie Parameters", template.CookieParams)
			}

			for statusCodeName, responseRef := range operation.Responses {
				statusCode, err := strconv.Atoi(statusCodeName)
				if err != nil {
					continue
				}

				for mediaTypeName, mediaType := range responseRef.Value.Content {
					for exampleName, exampleRef := range mediaType.Examples {
						template := container[operation.OperationID][exampleName]

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
