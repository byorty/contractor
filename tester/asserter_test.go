package tester_test

import (
	"github.com/byorty/contractor/tester"
	"github.com/stretchr/testify/suite"
	"testing"
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
	s.builder = tester.NewFxAsserterBuilder()
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

func (s *AsserterTestSuite) TestAsserter() {
	processor := s.builder.Build(s.testCase)
	processor.Process(s.testCase)
	s.T().Log(s.testCase.GetAssertions())
	s.Equal(tester.TestCaseStatusFailure, s.testCase.Status)
	//assertions := s.asserter.Assert(s.testCase)
	//s.T().Log(assertions)
}
