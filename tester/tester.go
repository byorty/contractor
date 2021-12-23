package tester

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/byorty/contractor/common"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	headerContentType = "Content-Type"
)

type Tester interface {
	Configure(ctx context.Context, containers common.TemplateContainer)
	Test() error
}

func NewFxTester() Tester {
	return &tester{
		suites: make([]TestSuite, 0),
		unmarshalers: map[string]func(buf []byte) (map[string]interface{}, error){
			"application/json": func(buf []byte) (map[string]interface{}, error) {
				obj := make(map[string]interface{})
				err := json.Unmarshal(buf, &obj)
				if err != nil {
					return nil, err
				}

				return obj, nil
			},
			"application/xml": func(buf []byte) (map[string]interface{}, error) {
				obj := make(map[string]interface{})
				err := xml.Unmarshal(buf, &obj)
				if err != nil {
					return nil, err
				}

				return obj, nil
			},
		},
	}
}

type tester struct {
	suites       []TestSuite
	unmarshalers map[string]func(buf []byte) (map[string]interface{}, error)
}

func (t *tester) Configure(ctx context.Context, containers common.TemplateContainer) {
	for suiteName, container := range containers {
		suite := TestSuite{
			Name:      suiteName,
			TestCases: make([]TestCase, 0),
		}

		for caseName, template := range container {
			for statusCode, expectedResponseByContentType := range template.ExpectedResponses {
				for contentType, example := range expectedResponseByContentType {
					tc := TestCase{
						Name:   caseName,
						Status: TestCaseStatusUndefined,
						ExpectedResult: TestCaseResult{
							StatusCode: statusCode,
							Headers: map[string]string{
								headerContentType: contentType,
							},
							Body: example,
						},
					}

					req, err := http.NewRequest(template.Method, template.GetUrl(), nil)
					if err == nil {
						req.Header.Add(headerContentType, contentType)
						req.URL.RawQuery = template.GetQueryParams().Encode()
						tc.request = req
					} else {
						tc.Err = err
						tc.Status = TestCaseStatusFailure
					}

					suite.TestCases = append(suite.TestCases, tc)
				}
			}
		}

		t.suites = append(t.suites, suite)
	}
}

func (t *tester) Test() error {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	for _, suite := range t.suites {
		for _, testCase := range suite.TestCases {
			//if testCase.Status == TestCaseStatusFailure {
			//	continue
			//}

			actualResult, err := t.runTestCase(client, testCase)
			if err == nil {
				testCase.ActualResult = actualResult
			} else {
				testCase.Status = TestCaseStatusFailure
				testCase.Err = err
			}
		}
	}

	return nil
}

func (t *tester) runTestCase(client http.Client, testCase TestCase) (*TestCaseResult, error) {
	resp, err := client.Do(testCase.request)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var body interface{}
	unmarshaler, ok := t.unmarshalers[resp.Header.Get(headerContentType)]
	if ok {
		body, err = unmarshaler(buf)
		if err != nil {
			return nil, err
		}
	}

	actualResult := &TestCaseResult{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Body:       body,
	}

	for headerName, _ := range testCase.ExpectedResult.Headers {
		actualResult.Headers[headerName] = resp.Header.Get(headerName)
	}

	return actualResult, nil
}

type TestSuite struct {
	Name      string
	TestCases []TestCase
}

type TestCase struct {
	Name           string
	Status         TestCaseStatus
	Err            error
	request        *http.Request
	response       *http.Response
	ExpectedResult TestCaseResult
	ActualResult   *TestCaseResult
}

type TestCaseStatus int

const (
	TestCaseStatusUndefined TestCaseStatus = iota
	TestCaseStatusSuccess
	TestCaseStatusFailure
)

type TestCaseResult struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}
