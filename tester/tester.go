package tester

import (
	"bytes"
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	//"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	variableExpr = "${%s}"
)

type Tester interface {
	Configure(ctx context.Context, arguments common.Arguments, containers common.TemplateContainer)
	Test() (TestCaseContainer, error)
}

func NewFxTester(
	loggerFactory common.LoggerFactory,
	mediaConverter common.MediaConverter,
	builder AsserterBuilder,
	postProcessorFactory PostProcessorFactory,
) Tester {
	return &tester{
		logger:               loggerFactory.CreateCommonLogger().Named("tester"),
		cases:                make(TestCaseContainer, 0),
		builder:              builder,
		mediaConverter:       mediaConverter,
		postProcessorFactory: postProcessorFactory,
	}
}

type tester struct {
	logger               common.Logger
	cases                TestCaseContainer
	builder              AsserterBuilder
	mediaConverter       common.MediaConverter
	postProcessorFactory PostProcessorFactory
	variables            map[string]string
}

func (t *tester) Configure(ctx context.Context, arguments common.Arguments, container common.TemplateContainer) {
	for templateName, template := range container {
		if !template.ContainsTags(arguments.Tags) {
			continue
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

				t.cases = append(t.cases, tc)
			}
		}
	}

	sort.Sort(t.cases)
	t.variables = arguments.Variables
}

func (t *tester) createRequest(tc *TestCase) (*http.Request, error) {
	template := tc.Template
	mediaType := tc.ExpectedResult.Headers[common.HeaderContentType]
	buf, err := t.mediaConverter.Marshal(common.MediaType(mediaType), tc.Template.Bodies[mediaType])
	if err != nil {
		return nil, err
	}

	url := template.GetUrl()
	query := template.GetQueryParams().Encode()
	headerParams := make(map[string]string)
	for key, value := range t.variables {
		url = strings.ReplaceAll(url, fmt.Sprintf(variableExpr, key), value)
		query = strings.ReplaceAll(query, fmt.Sprintf(variableExpr, key), value)
		buf = bytes.ReplaceAll(buf, []byte(fmt.Sprintf(variableExpr, key)), []byte(value))

		for headerName, rawHeaderValue := range template.HeaderParams {
			headerValue := fmt.Sprint(rawHeaderValue)
			headerValue = strings.ReplaceAll(headerValue, fmt.Sprintf(variableExpr, key), value)
			headerParams[headerName] = headerValue
		}
	}

	req, err := http.NewRequest(template.Method, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query
	req.Header.Add(common.HeaderAccept, mediaType)
	for headerName, headerValue := range headerParams {
		req.Header.Add(headerName, headerValue)
	}

	//t.logger.Debug(req.URL.String())
	//t.logger.Debug(req.Header)
	//t.logger.Debug(string(buf))

	return req, nil
}

func (t *tester) Test() (TestCaseContainer, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	for _, testCase := range t.cases {
		req, err := t.createRequest(testCase)
		if err == nil {
			testCase.request = req
		} else {
			testCase.Err = err
		}

		t.runTestCase(client, testCase)
		processor := t.builder.Build(testCase)
		processor.Process(testCase)

		if testCase.Status == TestCaseStatusSuccess {
			for _, def := range testCase.Template.PostProcessors {
				postProcessor, err := t.postProcessorFactory.Create(def)
				if err != nil {
					continue
				}

				postProcessor.PostProcess(testCase)
			}
		}
	}

	return t.cases, nil
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
		Buf:        buf,
	}

	for headerName, _ := range testCase.ExpectedResult.Headers {
		actualResult.Headers[headerName] = resp.Header.Get(headerName)
	}

	return actualResult, nil
}

type TestCaseContainer []*TestCase

func (c TestCaseContainer) Len() int {
	return len(c)
}

func (c TestCaseContainer) Less(i, j int) bool {
	return c[i].Template.Priority > c[j].Template.Priority
}

func (c TestCaseContainer) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c TestCaseContainer) HasError() bool {
	for _, testCase := range c {
		if testCase.Status == TestCaseStatusFailure {
			return true
		}
	}

	return false
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
	Assertions     []*Assertion
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
	Buf        []byte
}
