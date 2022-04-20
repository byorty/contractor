package assertion

import (
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/spyzhov/ajson"
)

type jsonContains struct {
	logger            common.Logger
	dataCrawler       common.DataCrawler
	dataCrawlerOpts   []common.DataCrawlerOption
	expressionFactory common.ExpressionFactory
	expressions       map[string]string
}

func NewJsonContains(
	logger common.Logger,
	dataCrawler common.DataCrawler,
	expressionFactory common.ExpressionFactory,
	definition map[string]interface{},
) tester.Assertion2 {
	expressions := make(map[string]string)
	dataCrawler.Walk(
		definition,
		func(k string, v interface{}) {
			expressions[k] = fmt.Sprint(v)
		},
		common.WithPrefix("$"),
		common.WithJoinKeys(),
		common.WithSkipCollections(),
	)
	return &jsonContains{
		logger:            logger.Named("json_assertion"),
		expressionFactory: expressionFactory,
		expressions:       expressions,
	}
}

func (a *jsonContains) Assert(data interface{}) (tester.AssertionResultList, error) {
	root, err := ajson.Unmarshal(data.([]byte))
	list := tester.NewAssertionResultList()
	if err != nil {
		list.Add(tester.AssertionResult{
			Status:   tester.AssertionResultStatusFailure,
			Expected: "json present",
			Actual:   err.Error(),
		})

		return list, nil
	}

	for path, expression := range a.expressions {
		output, err := a.expressionFactory.Create(common.ExpressionTypeAsserter, expression)
		if err != nil {
			list.Add(tester.AssertionResult{
				Status: tester.AssertionResultStatusFailure,
				Actual: err.Error(),
			})
			continue
		}

		asserter := output.(tester.Asserter)
		nodes, err := root.JSONPath(path)
		if err != nil {
			list.Add(tester.AssertionResult{
				Status:   tester.AssertionResultStatusFailure,
				Expected: asserter.GetExpected(),
				Actual:   err.Error(),
			})
			continue
		}

		if len(nodes) == 0 {
			list.Add(tester.AssertionResult{
				Status:   tester.AssertionResultStatusFailure,
				Expected: asserter.GetExpected(),
				Actual:   "nil",
			})
			continue
		}

		value, _ := nodes[0].Value()
		err = asserter.Assert(value)
		result := tester.AssertionResult{
			Status:   tester.AssertionResultStatusFailure,
			Expected: asserter.GetExpected(),
			Actual:   asserter.GetActual(),
		}

		if err == nil {
			result.Status = tester.AssertionResultStatusSuccess
		}

		list.Add(result)
	}

	return list, nil
}
