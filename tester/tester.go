package tester

import (
	"bytes"
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"net/http"
	"time"
)

type Tester interface {
	Configure(ctx context.Context, arguments common.Arguments, containers common.TemplateContainer)
	Test() (TestSuiteContainer, error)
}

func NewFxTester(
	mediaConverter common.MediaConverter,
	builder AsserterBuilder,
) Tester {
	return &tester{
		suites:         make([]TestSuite, 0),
		builder:        builder,
		mediaConverter: mediaConverter,
	}
}

type tester struct {
	suites         TestSuiteContainer
	builder        AsserterBuilder
	mediaConverter common.MediaConverter
}

func (t *tester) Configure(ctx context.Context, arguments common.Arguments, container common.TemplateContainer) {
	for templateName, template := range container {
		if !template.ContainsTags(arguments.Tags) {
			continue
		}

		suite := TestSuite{
			Name:      template.UID,
			TestCases: make([]*TestCase, 0),
		}

		for statusCode, expectedResponseByMediaType := range template.ExpectedResponses {
			for mediaTypeName, example := range expectedResponseByMediaType {
				tc := &TestCase{
					Name:     templateName,
					Status:   TestCaseStatusUndefined,
					Template: template,
					ExpectedResult: TestCaseResult{
						StatusCode: statusCode,
						Headers: map[string]string{
							common.HeaderContentType: mediaTypeName,
						},
						Body: example,
					},
				}

				req, err := t.createRequest(mediaTypeName, template)
				if err == nil {
					tc.request = req
				} else {
					tc.Err = err
				}

				suite.TestCases = append(suite.TestCases, tc)
			}
		}

		t.suites = append(t.suites, suite)
	}
}

func (t *tester) createRequest(mediaTypeName string, template *common.Template) (*http.Request, error) {
	buf, err := t.mediaConverter.Marshal(common.MediaType(mediaTypeName), template.Bodies[mediaTypeName])
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(template.Method, template.GetUrl(), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add(common.HeaderContentType, mediaTypeName)
	for headerName, headerValue := range template.HeaderParams {
		req.Header.Add(headerName, fmt.Sprint(headerValue))
	}

	req.URL.RawQuery = template.GetQueryParams().Encode()
	return req, nil
}

func (t *tester) Test() (TestSuiteContainer, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	for _, suite := range t.suites {
		for _, testCase := range suite.TestCases {
			t.runTestCase(client, testCase)
			processor := t.builder.Build(testCase)
			processor.Process(testCase)
		}
	}

	return t.suites, nil
}

func (t *tester) runTestCase(client *http.Client, testCase *TestCase) {
	if testCase.Err != nil {
		return
	}

	actualResult, err := t.sendRequest(client, testCase)
	if err == nil {
		testCase.ActualResult = actualResult
	} else {
		testCase.Err = err
	}
}

func (t *tester) sendRequest(client *http.Client, testCase *TestCase) (*TestCaseResult, error) {
	resp, err := client.Do(testCase.request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	body, err := t.mediaConverter.Unmarshal(common.MediaType(resp.Header.Get(common.HeaderContentType)), buf)
	if err != nil {
		return nil, err
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

type TestSuiteContainer []TestSuite

func (c TestSuiteContainer) HasError() bool {
	for _, suite := range c {
		for _, testCase := range suite.TestCases {
			if testCase.Status == TestCaseStatusFailure {
				return true
			}
		}
	}

	return false
}

type TestSuite struct {
	Name      string
	TestCases []*TestCase
}

type TestCase struct {
	Name           string
	Status         TestCaseStatus
	Err            error
	request        *http.Request
	response       *http.Response
	Template       *common.Template
	ExpectedResult TestCaseResult
	ActualResult   *TestCaseResult
	path           cmp.Path
	assertions     []*Assertion
}

func (t *TestCase) GetAssertions() []*Assertion {
	return t.assertions
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
