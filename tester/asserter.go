package tester

import (
	"fmt"
	"github.com/antonmedv/expr"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
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

func NewFxAsserterBuilder() AsserterBuilder {
	return &asserterBuilder{}
}

type asserterBuilder struct {
}

func (b *asserterBuilder) Build(testCase *TestCase) AssertionProcessor {
	assertions := AssertionMap{
		"__Status__": &Assertion{
			Name:       "Status Code",
			expression: fmt.Sprintf("NewEq(%d)", testCase.ExpectedResult.StatusCode),
			value:      testCase.ActualResult.StatusCode,
		},
		"__Err__": &Assertion{
			Name:       "Error",
			expression: "NewEmpty()",
			value:      testCase.Err,
		},
	}

	for expectedHeaderName, expectedHeaderValue := range testCase.ExpectedResult.Headers {
		assertions[fmt.Sprintf("__%s__", expectedHeaderName)] = &Assertion{
			Name:       fmt.Sprintf("Header %s", expectedHeaderName),
			expression: fmt.Sprintf("NewEq('%s')", expectedHeaderValue),
			value:      testCase.ActualResult.Headers[expectedHeaderName],
		}
	}

	b.buildMap(assertions, "", testCase.ExpectedResult.Body, func(m map[string]*Assertion, k string, v interface{}) {
		strVal := v.(string)
		m[k] = &Assertion{
			expression: fmt.Sprintf("New%s%s", strings.Title(strVal[0:1]), strVal[1:]),
		}
	})
	b.buildMap(assertions, "", testCase.ActualResult.Body, func(m map[string]*Assertion, k string, v interface{}) {
		_, ok := m[k]
		if !ok {
			m[k] = &Assertion{
				value: v,
			}
			return
		}

		m[k].value = v
	})

	return &assertionProcessor{
		assertions: assertions,
		env: map[string]interface{}{
			"NewEq": func(expected interface{}) Asserter {
				return &eqAsserter{
					expected: expected,
				}
			},
			"NewPositive": func() Asserter {
				return &minAsserter{
					expected: 1,
				}
			},
			"NewMin": func(expected float64) Asserter {
				return &minAsserter{
					expected: expected,
				}
			},
			"NewRegex": func(expr string) Asserter {
				return &regexAsserter{
					expected: expr,
				}
			},
			"NewEmpty": func() Asserter {
				return &emptyAsserter{}
			},
			"NewDate": func(layout string) Asserter {
				var expected string
				switch layout {
				case "RFC3339":
					expected = time.RFC3339
				case "RFC3339NANO":
					expected = time.RFC3339Nano
				default:
					expected = layout
				}
				return &dateAsserter{
					expected: expected,
				}
			},
		},
	}
}

func (b *asserterBuilder) buildMap(assertions AssertionMap, parentPath string, rawBody interface{}, fillFunc func(m map[string]*Assertion, k string, v interface{})) {
	switch body := rawBody.(type) {
	case map[string]interface{}:
		for k, val := range body {
			switch v := val.(type) {
			case map[string]interface{}, []interface{}:
				b.buildMap(
					assertions,
					fmt.Sprintf("%s.%s", parentPath, k),
					v,
					fillFunc,
				)
			default:
				fillFunc(assertions, fmt.Sprintf("%s.%s", parentPath, k), v)
			}
		}
	case []interface{}:
		for i, item := range body {
			b.buildMap(
				assertions,
				fmt.Sprintf("%s.%d", parentPath, i),
				item,
				fillFunc,
			)
		}
	}
}

type AssertionProcessor interface {
	Process(testCase *TestCase)
}

type assertionProcessor struct {
	assertions AssertionMap
	env        map[string]interface{}
}

func (a *assertionProcessor) Process(testCase *TestCase) {
	testCase.assertions = make([]*Assertion, 0)

	defer func() {
		if len(testCase.assertions) == 0 {
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
			testCase.assertions = append(testCase.assertions, assertion)
			continue
		}

		output, err := expr.Eval(assertion.expression, a.env)
		if err != nil {
			assertion.Actual = err.Error()
			testCase.assertions = append(testCase.assertions, assertion)
			continue
		}

		asserter := output.(Asserter)
		err = asserter.Assert(assertion.value)
		if err != nil {
			assertion.Expected = asserter.GetExpected()
			assertion.Actual = asserter.GetActual()
			testCase.assertions = append(testCase.assertions, assertion)
		}
	}
}

type Asserter interface {
	Assert(value interface{}) error
	GetExpected() string
	GetActual() string
}

type eqAsserter struct {
	expected interface{}
	actual   interface{}
}

func (a *eqAsserter) Assert(value interface{}) error {
	a.actual = value
	return validation.In(fmt.Sprint(a.expected)).Validate(fmt.Sprint(value))
}

func (a *eqAsserter) GetExpected() string {
	return fmt.Sprint(a.expected)
}

func (a *eqAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
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
	return fmt.Sprintf("great than %d", a.expected)
}

func (a *minAsserter) GetActual() string {
	return fmt.Sprint(a.actual)
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
