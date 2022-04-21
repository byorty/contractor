package graylog

import (
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/client"
	"github.com/byorty/contractor/tester/client/graylog/saved"
	"github.com/go-openapi/runtime"
)

func NewRunner(
	logger common.Logger,
	graylogClient client.GraylogClient,
	authInfo runtime.ClientAuthInfoWriter,
) tester.Runner {
	return &runner{
		logger:        logger.Named("graylog_test_runner"),
		graylogClient: graylogClient,
		authInfo:      authInfo,
	}
}

type runner struct {
	logger        common.Logger
	graylogClient client.GraylogClient
	request       *saved.SearchRelativeParams
	authInfo      runtime.ClientAuthInfoWriter
	testCase      tester.TestCase2
}

func (r *runner) Setup(testCase tester.TestCase2) {
	r.testCase = testCase
	sorting := "timestamp:asc"
	r.request = &saved.SearchRelativeParams{
		Query: fmt.Sprintf("%s", testCase.Setup.Query),
		Range: int64(testCase.Setup.Range.Seconds()),
		Sort:  &sorting,
	}
}

func (r *runner) Run(assertions tester.Asserter2List) tester.RunnerReportList {
	reports := tester.NewRunnerReportList()
	resp, err := r.graylogClient.SearchRelative(r.request, r.authInfo)
	if err != nil {
		r.logger.Error(err)
		assertionsResults := tester.NewAssertionResultList()
		assertionsResults.Add(tester.AssertionResult{
			Status:   tester.AssertionResultStatusFailure,
			Expected: "graylog response present",
			Actual:   err.Error(),
		})
		reports.Add(tester.RunnerReport{
			Name:       r.testCase.Name,
			Assertions: assertionsResults,
		})
		return reports
	}

	correlationMessages := common.NewMap[string, common.List[string]]()
	for _, item := range resp.Payload.Messages {
		message := item.Message.(map[string]interface{})
		correlationId := fmt.Sprint(message["correlation_id"])
		messages, ok := correlationMessages.Get(correlationId)
		if !ok {
			messages = common.NewList[string]()
			correlationMessages.Set(correlationId, messages)
		}

		messages.Add(fmt.Sprint(message["msg"]))
	}

	assertionIterator := assertions.Iterator()
	for correlationId, list := range correlationMessages.Entries() {
		var testCaseStarted bool
		assertionIterator.Reset()
		assertion := assertionIterator.Next()
		messageIterator := list.Iterator()
		correlationResults := tester.NewAssertionResultList()
		for messageIterator.HasNext() {
			message := messageIterator.Next()
			if !testCaseStarted && message == "start" {
				testCaseStarted = true
				continue
			}

			if !testCaseStarted {
				continue
			}

			results := assertion.Assert(message)
			correlationResults.Add(results.Entries()...)
			if assertionIterator.HasNext() && results.IsPassed() {
				assertion = assertionIterator.Next()
				continue
			}

			break
		}

		if correlationResults.IsPassed() {
			reports.RemoveAll()
			reports.Add(tester.RunnerReport{
				Name:       r.testCase.Name,
				Assertions: correlationResults,
			})
			return reports
		}

		if correlationResults.Len() == 0 {
			correlationResults.Add(tester.AssertionResult{
				Status:   tester.AssertionResultStatusFailure,
				Expected: "runner messages present",
				Actual:   "nil",
			})
		}

		reports.Add(tester.RunnerReport{
			Name:       fmt.Sprintf("%s#%s", r.testCase.Name, correlationId),
			Assertions: correlationResults,
		})
	}

	if reports.Len() == 0 {
		reports.Add(tester.RunnerReport{
			Name: r.testCase.Name,
			Assertions: tester.NewAssertionResultList(tester.AssertionResult{
				Status:   tester.AssertionResultStatusFailure,
				Expected: "graylog messages present",
				Actual:   "nil",
			}),
		})
	}

	return reports
}
