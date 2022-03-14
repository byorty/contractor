package tester_test

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/byorty/contractor/tester"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestAsserterTestSuite(t *testing.T) {
	suite.Run(t, new(AsserterTestSuite))
}

type AsserterTestSuite struct {
	suite.Suite
	builder  tester.AsserterBuilder
	testCase *tester.TestCase
}

func (s *AsserterTestSuite) SetupSuite() {
	s.testCase = &tester.TestCase{
		ExpectedResult: tester.TestCaseResult{
			StatusCode: 200,
			Body: map[string]interface{}{
				"str_eq": "eq('1234')",
				"int_eq": "eq(1234)",
				"map": map[string]interface{}{
					"key": "eq('value')",
				},
			},
		},
		ActualResult: &tester.TestCaseResult{
			StatusCode: 200,
			Body: map[string]interface{}{
				"str_eq": "12345",
				"int_eq": 12345,
				"map": map[string]interface{}{
					"key": "value1",
				},
			},
		},
	}
}

func (s *AsserterTestSuite) TestEqAsserter() {
	expected := randomdata.Alphanumeric(32)
	asserter := tester.NewEqAsserter(expected)
	s.NotNil(asserter.Assert(randomdata.Alphanumeric(16)))
	s.Nil(asserter.Assert(expected))
	s.Equal(expected, asserter.GetExpected())
	s.Equal(expected, asserter.GetActual())

	expectedInt := randomdata.Number(1, 10)
	asserter = tester.NewEqAsserter(expectedInt)
	s.NotNil(asserter.Assert(randomdata.Alphanumeric(16)))
	s.Nil(asserter.Assert(expectedInt))
	s.Nil(asserter.Assert(fmt.Sprint(expectedInt)))
	s.Nil(asserter.Assert(float64(expectedInt)))

	expectedFloat := randomdata.Decimal(1, 10)
	asserter = tester.NewEqAsserter(expectedFloat)
	s.Nil(asserter.Assert(expectedFloat))
}

func (s *AsserterTestSuite) TestPositiveAsserter() {
	asserter := tester.NewPositiveAsserter()
	actual := randomdata.Decimal(-100, 0)
	s.NotNil(asserter.Assert(actual))
	s.Nil(asserter.Assert(float64(1)))
	s.Nil(asserter.Assert(randomdata.Decimal(10, 1000)))
}

func (s *AsserterTestSuite) TestMinAsserter() {
	expected := randomdata.Decimal(100, 200)
	actual := randomdata.Decimal(1, 99)
	asserter := tester.NewMinAsserter(expected)
	s.NotNil(asserter.Assert(actual))
	s.Equal(fmt.Sprintf("great than %f", expected), asserter.GetExpected())
	s.Equal(fmt.Sprint(actual), asserter.GetActual())
	s.Nil(asserter.Assert(expected + 1))
	s.Nil(asserter.Assert(randomdata.Decimal(300, 999)))
}

func (s *AsserterTestSuite) TestMaxAsserter() {
	expected := randomdata.Decimal(100, 200)
	actual := randomdata.Decimal(300, 999)
	asserter := tester.NewMaxAsserter(expected)
	s.NotNil(asserter.Assert(actual))
	s.Equal(fmt.Sprintf("less than %f", expected), asserter.GetExpected())
	s.Equal(fmt.Sprint(actual), asserter.GetActual())
	s.Nil(asserter.Assert(expected - 1))
	s.Nil(asserter.Assert(randomdata.Decimal(1, 99)))
}

func (s *AsserterTestSuite) TestRegexAsserter() {
	expected := "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	actual := randomdata.Email()
	asserter := tester.NewRegexAsserter(expected)
	s.NotNil(asserter.Assert(randomdata.Alphanumeric(32)))
	s.Nil(asserter.Assert(actual))
	s.Equal(expected, asserter.GetExpected())
	s.Equal(actual, asserter.GetActual())
}

func (s *AsserterTestSuite) TestEmptyAsserter() {
	actual := randomdata.Email()
	asserter := tester.NewEmptyAsserter()
	s.NotNil(asserter.Assert(actual))
	s.Equal("nil", asserter.GetExpected())
	s.Equal(actual, asserter.GetActual())
	s.Nil(asserter.Assert(nil))
	s.Nil(asserter.Assert(map[string]interface{}{}))
	s.Nil(asserter.Assert(nil))
}

func (s *AsserterTestSuite) TestDateAsserter() {
	actual := randomdata.Email()
	asserter := tester.NewDateAsserter("RFC3339")
	s.NotNil(asserter.Assert(actual))
	s.Equal(fmt.Sprintf("format %s", time.RFC3339), asserter.GetExpected())
	s.Equal(actual, asserter.GetActual())
	s.Nil(asserter.Assert(time.Now().Format(time.RFC3339)))

	asserter = tester.NewDateAsserter("RFC3339NANO")
	s.Nil(asserter.Assert(time.Now().Format(time.RFC3339Nano)))

	asserter = tester.NewDateAsserter(time.Layout)
	s.Nil(asserter.Assert(time.Now().Format(time.Layout)))
}

func (s *AsserterTestSuite) TestRangeAsserter() {
	actual := randomdata.Number(1, 100)
	asserter := tester.NewRangeAsserter(1, 100)
	s.NotNil(asserter.Assert(randomdata.Number(-100, 0)))
	s.NotNil(asserter.Assert(randomdata.Number(101, 1000)))
	s.Equal(fmt.Sprintf("great than %v and less than %v", 1, 100), asserter.GetExpected())
	s.Nil(asserter.Assert(actual))
	s.Equal(fmt.Sprint(actual), asserter.GetActual())
}

func (s *AsserterTestSuite) TestContainsAsserter() {
	expected := randomdata.Alphanumeric(100)
	asserter := tester.NewContainsAsserter(expected)
	s.NotNil(asserter.Assert(randomdata.Alphanumeric(100)))
	s.Equal(fmt.Sprintf("must contains %s", expected), asserter.GetExpected())
	s.Nil(asserter.Assert(expected))
	s.Equal(expected, asserter.GetActual())
}
