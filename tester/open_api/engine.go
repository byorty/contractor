package open_api

import (
	"context"
	"encoding/json"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/config"
	"path/filepath"
	"strconv"
	"strings"
)

var _ tester.Engine = (*engine)(nil)

type engine struct {
	ctx                   context.Context
	logger                common.Logger
	configProviderFactory common.ConfigProviderFactory
	cfg                   *Config
	testCases             tester.TestCase2List
}

func (e *engine) Configure(data interface{}) error {
	configProvider, err := e.configProviderFactory.CreateByOptions(config.Static(data))
	if err != nil {
		e.logger.Error(err)
		return err
	}

	err = configProvider.Populate(e.cfg)
	if err != nil {
		e.logger.Error(err)
		return err
	}

	var spec *openapi3.T

	return nil
}

func (e *engine) processSpec(spec *openapi3.T) {
	for pathName, pathItem := range spec.Paths {
		for httpMethod, operation := range pathItem.Operations() {
			e.processXOperationExamples(operation)

			testCase := &tester.TestCase2{
				Name: "",
				//Tags:       nil,
				//Priority:   0,
				Setup: tester.TestCase2Setup{
					//Query:        "",
					//Range:        0,
					//Trigger:      "",
					//Headers:      nil,
					//Parameters:   nil,
					HeaderParams: make(map[string]interface{}),
					PathParams:   make(map[string]interface{}),
					QueryParams:  make(map[string]interface{}),
					CookieParams: make(map[string]interface{}),
				},
			}

			for _, parameterRef := range operation.Parameters {
				if parameterRef.Value == nil {
					continue
				}

				parameter := parameterRef.Value
				for _, exampleRef := range parameter.Examples {
					switch parameter.In {
					case openapi3.ParameterInPath:
						testCase.Setup.PathParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInQuery:
						testCase.Setup.QueryParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInHeader:
						testCase.Setup.HeaderParams[parameter.Name] = exampleRef.Value.Value
					case openapi3.ParameterInCookie:
						testCase.Setup.CookieParams[parameter.Name] = exampleRef.Value.Value
					}
				}
			}

			//if operation.RequestBody != nil {
			//	for mediaTypeName, mediaType := range operation.RequestBody.Value.Content {
			//		for exampleName, exampleRef := range mediaType.Examples {
			//			template := c.container.Create(exampleName, arguments.Url, pathName, httpMethod)
			//			template.Bodies[mediaTypeName] = exampleRef.Value.Value
			//		}
			//	}
			//}

			//for statusCodeName, responseRef := range operation.Responses {
			//	statusCode, err := strconv.Atoi(statusCodeName)
			//	if err != nil {
			//		continue
			//	}
			//
			//	for mediaTypeName, mediaType := range responseRef.Value.Content {
			//		if len(mediaType.Examples) == 0 {
			//			continue
			//		}
			//
			//		for exampleName, exampleRef := range mediaType.Examples {
			//			example := exampleRef.Value
			//			template := c.container[exampleName]
			//			if template == nil {
			//				continue
			//			}
			//
			//			_, ok := template.ExpectedResponses[statusCode]
			//			if !ok {
			//				template.ExpectedResponses[statusCode] = make(map[string]interface{})
			//			}
			//
			//			template.ExpectedResponses[statusCode][mediaTypeName] = example.Value
			//		}
			//	}
			//}
		}
	}
}

func (e *engine) processXOperationExamples(operation *openapi3.Operation) {
	examplesFilename, ok := operation.Extensions[xExamples]
	if !ok {
		return
	}

	var operationExamples tester.TestCase2
	exampleProvider, err := e.configProviderFactory.CreateByFiles(filepath.Join(
		filepath.Dir(e.cfg.Filename),
		strings.Trim(string(examplesFilename.(json.RawMessage)), "\""),
	))
	if err != nil {
		e.logger.Error(err)
		return
	}

	err = exampleProvider.Populate(&operationExamples)
	if err != nil {
		e.logger.Error(err)
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
						ExtensionProps: openapi3.ExtensionProps{
							Extensions: map[string]interface{}{
								xTagsName:       example.Tags,
								xPriority:       example.Priority,
								xPostProcessors: example.PostProcessors,
							},
						},
					},
				}
			}
		}
	}
}

func (e engine) GetTestCase2List() tester.TestCase2List {
	return e.testCases
}

func (e engine) CreateRunner() tester.Runner {
	//TODO implement me
	panic("implement me")
}
