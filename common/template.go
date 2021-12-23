package common

import (
	"fmt"
	"net/url"
	"strings"
)

type Template struct {
	StatusCode        int
	BaseUrl           string
	Path              string
	Method            string
	HeaderParams      map[string]interface{}
	PathParams        map[string]interface{}
	QueryParams       map[string]interface{}
	CookieParams      map[string]interface{}
	ExpectedResponses map[int]map[string]interface{}
}

func (t Template) GetUrl() string {
	u := fmt.Sprintf("%s%s", t.BaseUrl, t.Path)
	for paramName, paramValue := range t.PathParams {
		u = strings.ReplaceAll(u, fmt.Sprintf("{%s}", paramName), fmt.Sprint(paramValue))
	}

	return u
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

type TemplateContainer map[string]map[string]Template
