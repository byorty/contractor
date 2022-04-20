package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/client"
	"github.com/byorty/contractor/tester/client/graylog/saved"
	"github.com/ghodss/yaml"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"net/url"
	"path/filepath"
	//"github.com/pkg/errors"
	"io/ioutil"
)

type worker struct {
	logger           common.Logger
	assertionFactory tester.AssertionFactory
	graylogClient    client.GraylogClient
	testCases        tester.TestCaseContainer
}

func (w *worker) GetType() common.WorkerType {
	return common.WorkerTypeE2E
}

func (w *worker) Configure(ctx context.Context, arguments common.Arguments) error {
	files, err := ioutil.ReadDir(arguments.SpecLocation)
	if err != nil {
		w.logger.Error(err)
		//return errors.Wrap(err)
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src := make(tester.TestCaseContainer)
		err = w.readAndUnmarshal(file.Name(), &src)
		if err != nil {
			w.logger.Error(err)
			return err
		}

		for key, val := range src {
			w.testCases[key] = val
		}
	}

	return nil
}

func (w *worker) readAndUnmarshal(filename string, i interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	switch filepath.Ext(filename) {
	case ".json":
		err = json.Unmarshal(buf, i)
		if err != nil {
			return err
		}
	case ".yml", ".yaml":
		err = yaml.Unmarshal(buf, i)
		if err != nil {
			return err
		}
	default:

	}

	return nil
}

func (w *worker) Run() error {
	assertionResults := make([]tester.AssertionResult, 0)
	testRuns := make(TestRunContainer)
	for testCaseName, testCase := range w.testCases {
		var queryParams url.Values
		for key, value := range testCase.Setup.Query {
			queryParams.Set(key, fmt.Sprint(value))
		}

		testRuns[testCaseName] = &TestRun{
			Messages: common.NewMap[string, common.List[client.GraylogMessage]](),
		}
		resp, err := w.graylogClient.Saved.SearchRelative(&saved.SearchRelativeParams{
			//Decorate: nil,
			//Fields:   nil,
			//Filter:   nil,
			//Limit: nil,
			//Offset:   nil,
			Query: queryParams.Encode(),
			Range: int64(testCase.Setup.Range.Seconds()),
			//Sort:     nil,
			//Context:  nil,
		}, runtime.ClientAuthInfoWriterFunc(func(request runtime.ClientRequest, registry strfmt.Registry) error {
			return nil
		}))
		if err != nil {
			w.logger.Error(err)
			return err
		}

		for _, rawMessage := range resp.Payload.Messages {
			var graylogMessage client.GraylogMessage
			err := json.Unmarshal(rawMessage.Message.([]byte), &graylogMessage)
			if err != nil {
				w.logger.Error(err)
				return err
			}

			graylogMessages, ok := testRuns[testCaseName].Messages.Get(graylogMessage.CorrelationId)
			if !ok {
				graylogMessages = common.NewList[client.GraylogMessage]()
				testRuns[testCaseName].Messages.Set(graylogMessage.CorrelationId, graylogMessages)
			}

			graylogMessages.Add(graylogMessage)
		}

		var testCasePassed bool
		testRunResults := common.NewList[tester.AssertionResult]()
		assertions := common.NewList[tester.Assertion2]()
		assertionIterator := assertions.Iterator()
		correlationMessageIterator := testRuns[testCaseName].Messages.Iterator()
		for correlationMessageIterator.HasNext() && !testCasePassed {
			assertionIterator.Reset()
			messageIterator := correlationMessageIterator.Next().Iterator()

			var testCaseStarted bool
			correlationResults := tester.AssertionResultList(common.NewList[tester.AssertionResult]())
			assertion := assertionIterator.Next()
			for messageIterator.HasNext() {
				message := messageIterator.Next()
				if !testCaseStarted && message.Target == "test_case_start" {
					testCaseStarted = true
				}

				if !testCaseStarted {
					continue
				}

				results, _ := assertion.Assert([]byte(message.Msg))
				correlationResults.Add(results.Entries()...)
				if assertionIterator.HasNext() && results.IsPassed() {
					assertion = assertionIterator.Next()
				}
			}

			testCasePassed = !assertionIterator.HasNext() && correlationResults.IsPassed()
		}
	}
	return nil
}
