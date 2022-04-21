package reporter

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
)

func NewConsole(
	loggerFactory common.LoggerFactory,
) tester.Reporter2 {
	return &console{
		commonLogger:  loggerFactory.CreateCommonLogger(),
		successLogger: loggerFactory.CreateSuccessLogger(),
		errorLogger:   loggerFactory.CreateErrorLogger(),
	}
}

type console struct {
	commonLogger  common.Logger
	successLogger common.Logger
	errorLogger   common.Logger
}

func (r *console) Report(report tester.RunnerReport) {
	r.commonLogger.PrintGroup("Test Case: %s", report.Name)
	if report.Details != nil {
		for _, detail := range report.Details {
			switch data := detail.Data.(type) {
			case map[string]interface{}:
				r.commonLogger.PrintParameters(detail.Name, data)
			default:
				r.commonLogger.PrintParameter(detail.Name, data)
			}
		}
	}

	if report.Assertions.IsPassed() {
		r.successLogger.PrintParameter("Status", "Success")
	} else {
		r.errorLogger.PrintParameter("Status", "Failure")
		for _, assertion := range report.Assertions.Entries() {
			r.errorLogger.PrintParameters(assertion.Name, map[string]interface{}{
				common.LoggerReservedKeyExpected: assertion.Expected,
				common.LoggerReservedKeyActual:   assertion.Actual,
			})
		}
	}
}
