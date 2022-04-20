package runner

import (
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/client"
	"github.com/byorty/contractor/tester/client/graylog/saved"
	"github.com/go-openapi/runtime"
)

func NewFxGraylog(
	logger common.Logger,
	assertionFactory tester.AssertionFactory,
	graylogClient client.GraylogClient,
) tester.TestRunner {
	return &graylog{
		logger:           logger.Named("graylog_runner"),
		assertionFactory: assertionFactory,
		graylogClient:    graylogClient,
		authInfo:         nil,
	}
}

type graylog struct {
	logger           common.Logger
	assertionFactory tester.AssertionFactory
	graylogClient    client.GraylogClient
	request          *saved.SearchRelativeParams
	authInfo         runtime.ClientAuthInfoWriter
	name             string
}

func (r *graylog) Setup(name string, testCase tester.TestCaseDefinition) {
	r.name = name
	sorting := "timestamp:asc"
	r.request = &saved.SearchRelativeParams{
		Query: fmt.Sprintf("%s", testCase.Setup.Query),
		Range: int64(testCase.Setup.Range.Seconds()),
		Sort:  &sorting,
	}
}

func (r *graylog) Run(assertions tester.Assertion2List) tester.TestRunnerReportList {
	reports := tester.NewTestRunnerReportList()
	resp, err := r.graylogClient.SearchRelative(r.request, r.authInfo)
	if err != nil {
		r.logger.Error(err)
		assertionsResults := tester.NewAssertionResultList()
		assertionsResults.Add(tester.AssertionResult{
			Status:   tester.AssertionResultStatusFailure,
			Expected: "graylog response present",
			Actual:   err.Error(),
		})
		reports.Add(tester.TestRunnerReport{
			Name:       r.name,
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

			results, _ := assertion.Assert(message)
			correlationResults.Add(results.Entries()...)
			if assertionIterator.HasNext() && results.IsPassed() {
				assertion = assertionIterator.Next()
				continue
			}

			break
		}

		if correlationResults.IsPassed() {
			reports.RemoveAll()
			reports.Add(tester.TestRunnerReport{
				Name:       r.name,
				Assertions: correlationResults,
			})
			return reports
		}

		if correlationResults.Len() == 0 {
			correlationResults.Add(tester.AssertionResult{
				Status:   tester.AssertionResultStatusFailure,
				Expected: "graylog messages present",
				Actual:   "nil",
			})
		}

		reports.Add(tester.TestRunnerReport{
			Name:       fmt.Sprintf("%s#%s", r.name, correlationId),
			Assertions: correlationResults,
		})
	}

	if reports.Len() == 0 {
		reports.Add(tester.TestRunnerReport{
			Name: r.name,
			Assertions: tester.NewAssertionResultList(tester.AssertionResult{
				Status:   tester.AssertionResultStatusFailure,
				Expected: "graylog messages present",
				Actual:   "nil",
			}),
		})
	}

	return reports
}
