package common

import (
	"fmt"
	"net/url"
	"strings"
)

type Template struct {
	UID               string
	BaseUrl           string
	Path              string
	Method            string
	HeaderParams      map[string]interface{}
	PathParams        map[string]interface{}
	QueryParams       map[string]interface{}
	CookieParams      map[string]interface{}
	ExpectedResponses map[int]map[string]interface{}
	Bodies            map[string]interface{}
}

func (t Template) GetUrl() string {
	return fmt.Sprintf("%s%s", t.BaseUrl, t.GetPath())
}

func (t Template) GetPath() string {
	path := t.Path
	for paramName, paramValue := range t.PathParams {
		path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", paramName), fmt.Sprint(paramValue))
	}
	return path
}

func (t Template) GetQueryParams() url.Values {
	values := make(url.Values)
	for queryName, queryParam := range t.QueryParams {
		queryParams, ok := queryParam.([]interface{})
		if ok {
			for _, queryParam := range queryParams {
				values.Add(queryName, fmt.Sprint(queryParam))
			}
		} else {
			values.Add(queryName, fmt.Sprint(queryParam))
		}
	}

	return values
}

type TemplateContainer map[string]*Template

func (c TemplateContainer) Create(exampleName, baseUrl, pathName, httpMethod string) *Template {
	template, ok := c[exampleName]
	if !ok {
		template = &Template{
			BaseUrl:           baseUrl,
			Path:              pathName,
			Method:            httpMethod,
			PathParams:        make(map[string]interface{}),
			QueryParams:       make(map[string]interface{}),
			HeaderParams:      make(map[string]interface{}),
			CookieParams:      make(map[string]interface{}),
			ExpectedResponses: make(map[int]map[string]interface{}),
			Bodies:            make(map[string]interface{}),
		}

		c[exampleName] = template
	}

	return template
}
