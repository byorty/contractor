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

const (
	headerContentType = "Content-Type"
)

type Tester interface {
	Configure(ctx context.Context, containers common.TemplateContainer)
	Test() (TestSuiteContainer, error)
}

func NewFxTester(
	mediaConverter common.MediaConverter,
) Tester {
	return &tester{
		suites:         make([]TestSuite, 0),
		mediaConverter: mediaConverter,
	}
}

type tester struct {
	suites         TestSuiteContainer
	mediaConverter common.MediaConverter
}

func (t *tester) Configure(ctx context.Context, container common.TemplateContainer) {
	for templateName, template := range container {
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
							headerContentType: mediaTypeName,
						},
						Body: example,
					},
				}

				req, err := t.createRequest(mediaTypeName, template, example)
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

func (t *tester) createRequest(mediaTypeName string, template common.Template, example interface{}) (*http.Request, error) {
	buf, err := t.mediaConverter.Marshal(common.MediaType(mediaTypeName), example)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(template.Method, template.GetUrl(), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add(headerContentType, mediaTypeName)
	req.URL.RawQuery = template.GetQueryParams().Encode()
	return req, nil
}

func (t *tester) Test() (TestSuiteContainer, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	for _, suite := range t.suites {
		for _, testCase := range suite.TestCases {
			t.runTestCase(client, testCase)
		}
	}

	return t.suites, nil
}

func (t *tester) runTestCase(client http.Client, testCase *TestCase) {
	defer testCase.Assert()
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

func (t *tester) sendRequest(client http.Client, testCase *TestCase) (*TestCaseResult, error) {
	resp, err := client.Do(testCase.request)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := t.mediaConverter.Unmarshal(common.MediaType(resp.Header.Get(headerContentType)), buf)
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
	Template       common.Template
	ExpectedResult TestCaseResult
	ActualResult   *TestCaseResult
	path           cmp.Path
	assertions     []TestCaseAssertion
}

func (t *TestCase) PushStep(ps cmp.PathStep) {
	t.path = append(t.path, ps)
}

func (t *TestCase) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := t.path.Last().Values()
		t.assertions = append(t.assertions, TestCaseAssertion{
			Name:     fmt.Sprintf("Property %s value is not equal", t.path.GoString()),
			Expected: fmt.Sprint(vx),
			Actual:   fmt.Sprint(vy),
		})
	}
}

func (t *TestCase) PopStep() {
	t.path = t.path[:len(t.path)-1]
}

func (t *TestCase) Assert() {
	t.assertions = make([]TestCaseAssertion, 0)
	defer func() {
		if len(t.assertions) == 0 {
			t.Status = TestCaseStatusSuccess
		} else {
			t.Status = TestCaseStatusFailure
		}
	}()
	if t.Err != nil {
		t.assertions = append(t.assertions, TestCaseAssertion{
			Name:     "Error",
			Expected: "nil",
			Actual:   t.Err.Error(),
		})

		return
	}

	if t.ExpectedResult.StatusCode != t.ActualResult.StatusCode {
		t.assertions = append(t.assertions, TestCaseAssertion{
			Name:     "Status Code",
			Expected: fmt.Sprint(t.ExpectedResult.StatusCode),
			Actual:   fmt.Sprint(t.ActualResult.StatusCode),
		})
	}

	for expectedHeaderName, expectedHeaderValue := range t.ExpectedResult.Headers {
		actualHeaderValue, ok := t.ActualResult.Headers[expectedHeaderName]
		if ok && expectedHeaderValue == actualHeaderValue {
			continue
		}

		t.assertions = append(t.assertions, TestCaseAssertion{
			Name:     fmt.Sprintf("Header %s", expectedHeaderName),
			Expected: expectedHeaderValue,
			Actual:   actualHeaderValue,
		})
	}

	cmp.Equal(t.ExpectedResult.Body, t.ActualResult.Body, cmp.Reporter(t))

	return
}

func (t *TestCase) GetAssertions() []TestCaseAssertion {
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

type TestCaseAssertion struct {
	Name     string
	Expected string
	Actual   string
}
