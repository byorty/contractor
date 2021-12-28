package tester

import (
	"fmt"
	"github.com/byorty/contractor/logger"
)

type Reporter interface {
	Report(container TestSuiteContainer) error
}

func NewFxStdoutReporter(
	l logger.Logger,
) Reporter {
	return &stdoutReporter{
		logger: l,
	}
}

type stdoutReporter struct {
	logger logger.Logger
}

func (r *stdoutReporter) Report(container TestSuiteContainer) error {
	errorLogger := r.logger.ToErrorColors()
	for _, testSuite := range container {
		for _, testCase := range testSuite.TestCases {
			r.logger.PrintGroup("Test Case: %s", testCase.Name)
			r.logger.PrintSubGroup("Status Code: %d", testCase.ExpectedResult.StatusCode)
			r.logger.PrintSubGroup("Path: %s", testCase.Template.Path)
			r.logger.PrintParameters("Header Parameters", testCase.Template.HeaderParams)
			r.logger.PrintParameters("Path Parameters", testCase.Template.PathParams)
			r.logger.PrintParameters("Query Parameters", testCase.Template.QueryParams)
			r.logger.PrintParameters("Cookie Parameters", testCase.Template.CookieParams)

			if testCase.Status == TestCaseStatusSuccess {
				r.logger.PrintSuccess()
			} else {
				assertions := testCase.GetAssertions()
				r.logger.PrintFailure()
				for _, assertion := range assertions {
					errorLogger.PrintSubGroup(fmt.Sprintf("%s:", assertion.Name))
					errorLogger.PrintParameter("Expected", assertion.Expected)
					errorLogger.PrintParameter("Actual", assertion.Actual)
				}
			}

			fmt.Println()
		}
	}

	return nil
}
