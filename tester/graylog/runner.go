package graylog

import (
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/graylog/client/saved"
	"github.com/go-openapi/runtime"
	"time"
)

const (
	runnerName = "Graylog"
)

type Message struct {
	Body      string
	Timestamp time.Time
}

func NewRunner(
	ctx context.Context,
	logger common.Logger,
	graylogClient Client,
	authInfo runtime.ClientAuthInfoWriter,
) tester.Runner {
	return &runner{
		ctx:           ctx,
		logger:        logger.Named("graylog_test_runner"),
		graylogClient: graylogClient,
		authInfo:      authInfo,
	}
}

type runner struct {
	ctx           context.Context
	logger        common.Logger
	graylogClient Client
	request       *saved.SearchRelativeParams
	authInfo      runtime.ClientAuthInfoWriter
	testCase      tester.TestCase2
}

func (r *runner) Setup(testCase tester.TestCase2) {
	r.testCase = testCase
	var limit int64 = 1000
	sorting := "timestamp:asc"
	r.request = &saved.SearchRelativeParams{
		Context: r.ctx,
		Query:   fmt.Sprintf("%s AND test_case:true", testCase.Setup.Query),
		Range:   int64(testCase.Setup.Range.Seconds()),
		Sort:    &sorting,
		Limit:   &limit,
	}
}

func (r *runner) Run(assertions tester.Asserter2List) tester.RunnerReportList {
	reports := tester.NewRunnerReportList()
	resp, err := r.graylogClient.SearchRelative(r.request, r.authInfo)
	if err != nil {
		r.logger.Error(err)
		assertionsResults := tester.NewAssertionResultList()
		assertionsResults.Add(tester.AssertionResult{
			Name:     runnerName,
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

	uniqMessages := common.NewMap[string, bool]()
	correlationMessages := common.NewMap[string, common.List[Message]]()
	for _, item := range resp.Payload.Messages {
		message := item.Message.(map[string]interface{})
		if message["correlation_id"] == nil {
			continue
		}

		messageId := fmt.Sprint(message["_id"])
		//r.logger.Debug(item.Message)
		_, ok := uniqMessages.Get(messageId)
		if ok {
			continue
		}

		uniqMessages.Set(messageId, true)
		//r.logger.Debug(message["correlation_id"], message["msg"])
		correlationId := fmt.Sprint(message["correlation_id"])
		messageBody := fmt.Sprint(message["msg"])
		messages, ok := correlationMessages.Get(correlationId)
		if !ok {
			assert := assertions.Get(0)
			if assert == nil {
				continue
			}

			results := assert.Assert(messageBody)
			if results.IsPassed() {
				messages = common.NewList[Message]()
				correlationMessages.Set(correlationId, messages)
			}

			continue
		}

		timestamp, _ := time.Parse(time.RFC3339Nano, fmt.Sprint(message["timestamp"]))
		messages.Add(Message{
			Body:      messageBody,
			Timestamp: timestamp,
		})
	}

	//r.logger.Debug(uniqMessages)

	assertionIterator := assertions.Iterator()
	for correlationId, list := range correlationMessages.Entries() {
		assertionIterator.Reset()
		assertion := assertionIterator.Next()
		list.Sort(func(a, b Message) bool {
			return a.Timestamp.UnixNano() < b.Timestamp.UnixNano()
		})
		messageIterator := list.Iterator()
		correlationResults := tester.NewAssertionResultList()
		for messageIterator.HasNext() {
			message := messageIterator.Next()
			results := assertion.Assert(message.Body)
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
				Name:     runnerName,
				Status:   tester.AssertionResultStatusFailure,
				Expected: "graylog messages present",
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
				Name:     runnerName,
				Status:   tester.AssertionResultStatusFailure,
				Expected: "graylog messages present",
				Actual:   "nil",
			}),
		})
	}

	return reports
}
