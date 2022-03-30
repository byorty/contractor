package tester

import (
	"fmt"
	"github.com/byorty/contractor/common"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Assertion struct {
	Name       string
	Expected   string
	Actual     string
	expression string
	value      interface{}
}

type AssertionMap map[string]*Assertion

type AsserterBuilder interface {
	Build(testCase *TestCase) AssertionProcessor
}

func NewFxAsserterBuilder(
	dataCrawler common.DataCrawler,
	expressionFactory common.ExpressionFactory,
) AsserterBuilder {
	return &asserterBuilder{
		expressionFactory: expressionFactory,
		dataCrawler:       dataCrawler,
		dataCrawlerOpts: []common.DataCrawlerOption{
			common.WithJoinKeys(),
			common.WithSkipCollections(),
		},
	}
}

type asserterBuilder struct {
	dataCrawler       common.DataCrawler
	dataCrawlerOpts   []common.DataCrawlerOption
	expressionFactory common.ExpressionFactory
}

func (b *asserterBuilder) Build(testCase *TestCase) AssertionProcessor {
	processor := &assertionProcessor{
		expressionFactory: b.expressionFactory,
		assertions: AssertionMap{
			"__Err__": &Assertion{
				Name:       "Error",
				expression: "empty()",
				value:      testCase.Err,
			},
		},
	}

	if testCase.ActualResult == nil {
		return processor
	}

	processor.assertions["__Status__"] = &Assertion{
		Name:       "Status Code",
		expression: fmt.Sprintf("eq(%d)", testCase.ExpectedResult.StatusCode),
		value:      testCase.ActualResult.StatusCode,
	}

	for expectedHeaderName, expectedHeaderValue := range testCase.ExpectedResult.Headers {
		processor.assertions[fmt.Sprintf("__%s__", expectedHeaderName)] = &Assertion{
			Name:       fmt.Sprintf("Header %s", expectedHeaderName),
			expression: fmt.Sprintf("eq('%s')", expectedHeaderValue),
			value:      testCase.ActualResult.Headers[expectedHeaderName],
		}
	}

	b.dataCrawler.Walk(testCase.ExpectedResult.Body, func(k string, v interface{}) {
		processor.assertions[k] = &Assertion{
			expression: v.(string),
		}
	}, b.dataCrawlerOpts...)

	b.dataCrawler.Walk(testCase.ActualResult.Body, func(k string, v interface{}) {
		_, ok := processor.assertions[k]
		if !ok {
			processor.assertions[k] = &Assertion{
				value: v,
			}
			return
		}

		processor.assertions[k].value = v
	}, b.dataCrawlerOpts...)

	return processor
}

type AssertionProcessor interface {
	Process(testCase *TestCase)
}

type assertionProcessor struct {
	expressionFactory common.ExpressionFactory
	assertions        AssertionMap
}

func (a *assertionProcessor) Process(testCase *TestCase) {
	testCase.Assertions = make([]*Assertion, 0)

	defer func() {
		if len(testCase.Assertions) == 0 {
			testCase.Status = TestCaseStatusSuccess
		} else {
			testCase.Status = TestCaseStatusFailure
		}
	}()

	for path, assertion := range a.assertions {
		if len(assertion.Name) == 0 {
			assertion.Name = fmt.Sprintf("Property '%s' value is invalid", strings.TrimLeft(path, "."))
		}

		if len(assertion.expression) == 0 {
			assertion.Expected = "not present"
			assertion.Actual = "present"
			testCase.Assertions = append(testCase.Assertions, assertion)
			continue
		}

		output, err := a.expressionFactory.Create(common.ExpressionTypeAsserter, assertion.expression)
		if err != nil {
			assertion.Actual = err.Error()
			testCase.Assertions = append(testCase.Assertions, assertion)
			continue
		}

		asserter := output.(Asserter)
		err = asserter.Assert(assertion.value)
		if err != nil {
			assertion.Expected = asserter.GetExpected()
			assertion.Actual = asserter.GetActual()
			testCase.Assertions = append(testCase.Assertions, assertion)
		}
	}
}

type Asserter interface {
	Assert(value interface{}) error
	GetExpected() string
	GetActual() string
}

func NewEqAsserter(expected interface{}) Asserter {
	return &eqAsserter{
		expected: expected,
	}
}

type eqAsserter struct {
	expected interface{}
	actual   interface{}
}

func (a *eqAsserter) Assert(value interface{}) error {
	switch a.expected.(type) {
	case int8, int, int16, int32, int64:
		switch v := value.(type) {
		case float64:
			a.actual = int(v)
		case string:
			i, err := strconv.Atoi(fmt.Sprintf("%v", value))
			if err != nil {
				return err
			}

			a.actual = i
		default:
			a.actual = value
		}

	default:
		a.actual = value
	}

	return validation.In(a.expected).Validate(a.actual)
}

func (a *eqAsserter) GetExpected() string {
	return fmt.Sprint(a.expected)
}

func (a *eqAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewPositiveAsserter() Asserter {
	return &minAsserter{
		expected: 1,
	}
}

func NewMinAsserter(expected float64) Asserter {
	return &minAsserter{
		expected: expected,
	}
}

type minAsserter struct {
	expected float64
	actual   interface{}
}

func (a *minAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.Min(a.expected).Validate(value)
}

func (a *minAsserter) GetExpected() string {
	return fmt.Sprintf("great than %f", a.expected)
}

func (a *minAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewMaxAsserter(expected float64) Asserter {
	return &maxAsserter{
		expected: expected,
	}
}

type maxAsserter struct {
	expected float64
	actual   interface{}
}

func (a *maxAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.Max(a.expected).Validate(value)
}

func (a *maxAsserter) GetExpected() string {
	return fmt.Sprintf("less than %f", a.expected)
}

func (a *maxAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewRangeAsserter(min, max interface{}) Asserter {
	return &rangeAsserter{
		min: min,
		max: max,
	}
}

type rangeAsserter struct {
	min    interface{}
	max    interface{}
	actual interface{}
}

func (a *rangeAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.Validate(value, validation.Min(a.min), validation.Max(a.max))
}

func (a *rangeAsserter) GetExpected() string {
	return fmt.Sprintf("great than %v and less than %v", a.min, a.max)
}

func (a *rangeAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewRegexAsserter(expr string) Asserter {
	return &regexAsserter{
		expected: expr,
	}
}

type regexAsserter struct {
	expected string
	actual   interface{}
}

func (a *regexAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.Match(regexp.MustCompile(a.expected)).Validate(value)
}

func (a *regexAsserter) GetExpected() string {
	return a.expected
}

func (a *regexAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewEmptyAsserter() Asserter {
	return &emptyAsserter{}
}

type emptyAsserter struct {
	actual interface{}
}

func (a *emptyAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.Empty.Validate(value)
}

func (a *emptyAsserter) GetExpected() string {
	return "nil"
}

func (a *emptyAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewDateAsserter(layout string) Asserter {
	var expected string
	switch layout {
	case common.ExpressionDateLayoutRFC3339:
		expected = time.RFC3339
	case common.ExpressionDateLayoutRFC3339Nano:
		expected = time.RFC3339Nano
	default:
		expected = layout
	}

	return &dateAsserter{
		expected: expected,
	}
}

type dateAsserter struct {
	expected string
	actual   interface{}
}

func (a *dateAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.Date(a.expected).Validate(value)
}

func (a *dateAsserter) GetExpected() string {
	return fmt.Sprintf("format %s", a.expected)
}

func (a *dateAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}

func NewContainsAsserter(expr string) Asserter {
	return &containsAsserter{
		expected: expr,
	}
}

type containsAsserter struct {
	expected string
	actual   interface{}
}

func (a *containsAsserter) Assert(value interface{}) error {
	a.actual = value
	if strings.Contains(fmt.Sprint(value), a.expected) {
		return nil
	}

	return errors.New(a.GetExpected())
}

func (a *containsAsserter) GetExpected() string {
	return fmt.Sprintf("must contains %s", a.expected)
}

func (a *containsAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
}
