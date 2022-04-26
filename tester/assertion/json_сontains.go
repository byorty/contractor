package assertion

import (
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/spyzhov/ajson"
)

const (
	jsonContainsName = "Json Contains"
)

type jsonContains struct {
	logger            common.Logger
	dataCrawler       common.DataCrawler
	dataCrawlerOpts   []common.DataCrawlerOption
	expressionFactory common.ExpressionFactory
	name              string
	expressions       map[string]string
}

func NewJsonContains(
	logger common.Logger,
	dataCrawler common.DataCrawler,
	expressionFactory common.ExpressionFactory,
	name string,
	definition interface{},
) tester.Asserter2 {
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
		name:              name,
	}
}

func (a *jsonContains) Assert(data interface{}) tester.AssertionResultList {
	list := tester.NewAssertionResultList()
	root, err := ajson.Unmarshal([]byte(fmt.Sprint(data)))
	if err != nil {
		a.logger.Error(err)
		list.Add(tester.AssertionResult{
			Name:     a.name,
			Status:   tester.AssertionResultStatusFailure,
			Expected: "json present",
			Actual:   err.Error(),
		})

		return list
	}

	for path, expression := range a.expressions {
		resultName := fmt.Sprintf("Path '%s'", path)
		output, err := a.expressionFactory.Create(common.ExpressionTypeAsserter, expression)
		if err != nil {
			a.logger.Error(err)
			list.Add(tester.AssertionResult{
				Name:   resultName,
				Status: tester.AssertionResultStatusFailure,
				Actual: err.Error(),
			})
			continue
		}

		asserter := output.(tester.Asserter)
		nodes, err := root.JSONPath(path)
		if err != nil {
			a.logger.Error(err)
			list.Add(tester.AssertionResult{
				Name:     resultName,
				Status:   tester.AssertionResultStatusFailure,
				Expected: asserter.GetExpected(),
				Actual:   err.Error(),
			})
			continue
		}

		if len(nodes) == 0 {
			list.Add(tester.AssertionResult{
				Name:     resultName,
				Status:   tester.AssertionResultStatusFailure,
				Expected: asserter.GetExpected(),
				Actual:   "nil",
			})
			continue
		}

		for _, node := range nodes {
			value, _ := node.Value()
			err = asserter.Assert(value)
			result := tester.AssertionResult{
				Name:     resultName,
				Status:   tester.AssertionResultStatusSuccess,
				Expected: asserter.GetExpected(),
				Actual:   asserter.GetActual(),
			}

			if err != nil {
				a.logger.Error(err)
				result.Status = tester.AssertionResultStatusFailure
			}

			list.Add(result)
		}
	}

	return list
}
