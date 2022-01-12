package converter

import (
	"context"
	"encoding/json"
	"github.com/byorty/contractor/common"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type XOperationRequestExample struct {
	Headers    map[string]interface{} `json:"headers"`
	Parameters map[string]interface{} `json:"parameters"`
}

type XOperationResponseExample struct {
	StatusCode int                    `json:"status_code"`
	Body       map[string]interface{} `json:"body"`
}

type xOperationExample struct {
	Request  XOperationRequestExample  `json:"request"`
	Response XOperationResponseExample `json:"response"`
}

type xOperationExamples map[string]xOperationExample

type oaConverter struct {
	container common.TemplateContainer
}

func (c *oaConverter) convert(ctx context.Context, arguments common.Arguments, spec *openapi3.T) (common.TemplateContainer, error) {
	//err := spec.Validate(ctx)
	//if err != nil {
	//	return nil, err
	//}

	for pathName, pathItem := range spec.Paths {
		for httpMethod, operation := range pathItem.Operations() {
			c.processOperation(arguments, pathName, httpMethod, operation)
		}
	}

	return c.container, nil
}

func (c *oaConverter) processOperation(arguments common.Arguments, pathName, httpMethod string, operation *openapi3.Operation) {
	c.processXOperationExamples(arguments, operation)

	for _, parameterRef := range operation.Parameters {
		if parameterRef.Value == nil {
			continue
		}

		c.processParameter(arguments.BaseUrl, pathName, httpMethod, parameterRef.Value)
	}

	if operation.RequestBody != nil {
		for mediaTypeName, mediaType := range operation.RequestBody.Value.Content {
			for exampleName, exampleRef := range mediaType.Examples {
				template := c.container.Create(exampleName, arguments.BaseUrl, pathName, httpMethod)
				template.Bodies[mediaTypeName] = exampleRef.Value.Value
			}
		}
	}

	for statusCodeName, responseRef := range operation.Responses {
		statusCode, err := strconv.Atoi(statusCodeName)
		if err != nil {
			continue
		}

		for mediaTypeName, mediaType := range responseRef.Value.Content {
			if len(mediaType.Examples) == 0 {
				continue
			}

			for exampleName, exampleRef := range mediaType.Examples {
				template := c.container[exampleName]
				if template == nil {
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

func (c *oaConverter) processXOperationExamples(arguments common.Arguments, operation *openapi3.Operation) {
	examplesFilename, ok := operation.Extensions["x-examples"]
	if !ok {
		return
	}

	var operationExamples xOperationExamples
	err := c.readAndUnmarshal(
		filepath.Join(
			filepath.Dir(arguments.SpecFilename),
			strings.Trim(string(examplesFilename.(json.RawMessage)), "\""),
		),
		&operationExamples,
	)
	if err != nil {
		log.Println(err)
		return
	}

	for exampleName, example := range operationExamples {
		for headerName, headerValue := range example.Request.Headers {
			operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     headerName,
					In:       openapi3.ParameterInHeader,
					Required: true,
					Examples: openapi3.Examples{
						exampleName: &openapi3.ExampleRef{
							Value: &openapi3.Example{
								Value: headerValue,
							},
						},
					},
				},
			})
		}

		for paramName, paramValue := range example.Request.Parameters {
			if paramName == "body" && operation.RequestBody != nil {
				for _, mediaType := range operation.RequestBody.Value.Content {
					if mediaType.Examples == nil {
						mediaType.Examples = make(openapi3.Examples)
					}

					mediaType.Examples[exampleName] = &openapi3.ExampleRef{
						Value: &openapi3.Example{
							Value: paramValue,
						},
					}
				}
				continue
			}

			for _, parameterRef := range operation.Parameters {
				if paramName != parameterRef.Value.Name {
					continue
				}

				if parameterRef.Value.Examples == nil {
					parameterRef.Value.Examples = make(openapi3.Examples)
				}

				parameterRef.Value.Examples[exampleName] = &openapi3.ExampleRef{
					Value: &openapi3.Example{
						Value: paramValue,
					},
				}
			}
		}

		for statusCodeName, responseRef := range operation.Responses {
			if responseRef.Value == nil {
				continue
			}

			statusCode, err := strconv.Atoi(statusCodeName)
			if err != nil {
				continue
			}

			if example.Response.StatusCode != statusCode {
				continue
			}

			for _, mediaType := range responseRef.Value.Content {
				if mediaType.Examples == nil {
					mediaType.Examples = make(openapi3.Examples)
				}
				mediaType.Examples[exampleName] = &openapi3.ExampleRef{
					Value: &openapi3.Example{
						Value: example.Response.Body,
					},
				}
			}
		}
	}
}

func (c *oaConverter) processParameter(baseUrl, pathName, httpMethod string, parameter *openapi3.Parameter) {
	for exampleName, exampleRef := range parameter.Examples {
		template := c.container.Create(exampleName, baseUrl, pathName, httpMethod)

		switch parameter.In {
		case openapi3.ParameterInPath:
			template.PathParams[parameter.Name] = exampleRef.Value.Value
		case openapi3.ParameterInQuery:
			template.QueryParams[parameter.Name] = exampleRef.Value.Value
		case openapi3.ParameterInHeader:
			template.HeaderParams[parameter.Name] = exampleRef.Value.Value
		case openapi3.ParameterInCookie:
			template.CookieParams[parameter.Name] = exampleRef.Value.Value
		}
	}
}

func (c *oaConverter) readAndUnmarshal(filename string, i interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	switch filepath.Ext(filename) {
	case ".json":
		err = json.Unmarshal(buf, i)
		if err != nil {
			return err
		}
	case ".yml", ".yaml":
		err = yaml.Unmarshal(buf, i)
		if err != nil {
			return err
		}
	default:

	}

	return nil
}
