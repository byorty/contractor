package tester

import (
	"fmt"
	"github.com/byorty/contractor/common"
	"strings"
)

type Reporter interface {
	Report(container TestCaseContainer) error
}

func NewFxStdoutReporter(
	loggerFactory common.LoggerFactory,
) Reporter {
	return &stdoutReporter{
		commonLogger:  loggerFactory.CreateCommonLogger(),
		successLogger: loggerFactory.CreateSuccessLogger(),
		errorLogger:   loggerFactory.CreateErrorLogger(),
	}
}

type stdoutReporter struct {
	commonLogger  common.Logger
	successLogger common.Logger
	errorLogger   common.Logger
}

func (r *stdoutReporter) Report(container TestCaseContainer) error {
	for _, testCase := range container {
		r.commonLogger.PrintGroup("Test Case: %s", testCase.Name)
		r.commonLogger.PrintParameter("Path", testCase.Template.Path)
		r.commonLogger.PrintParameter("Method", strings.ToUpper(testCase.Template.Method))
		r.commonLogger.PrintParameter("Status Code", testCase.ExpectedResult.StatusCode)
		r.commonLogger.PrintParameters("Header Parameters", testCase.Template.HeaderParams)
		r.commonLogger.PrintParameters("Path Parameters", testCase.Template.PathParams)
		r.commonLogger.PrintParameters("Query Parameters", testCase.Template.QueryParams)
		r.commonLogger.PrintParameters("Cookie Parameters", testCase.Template.CookieParams)

		if testCase.Status == TestCaseStatusSuccess {
			r.successLogger.PrintParameter("Status", "Success")
		} else {
			r.errorLogger.PrintParameter("Status", "Failure")
			for _, assertion := range testCase.Assertions {
				r.errorLogger.PrintParameters(assertion.Name, map[string]interface{}{
					common.LoggerReservedKeyExpected: assertion.Expected,
					common.LoggerReservedKeyActual:   assertion.Actual,
				})
			}
		}

		fmt.Println()
	}

	return nil
}
