package assertion

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
)

func NewPlain(
	expressionFactory common.ExpressionFactory,
	name string,
	definition interface{},
) (tester.Asserter2, error) {
	out, err := expressionFactory.Create(common.ExpressionTypeAsserter, definition)
	if err != nil {
		return nil, err
	}

	return &plain{
		name:     name,
		asserter: out.(tester.Asserter),
	}, nil
}

type plain struct {
	name     string
	asserter tester.Asserter
}

func (a *plain) Assert(data interface{}) tester.AssertionResultList {
	status := tester.AssertionResultStatusFailure
	err := a.asserter.Assert(data)
	if err == nil {
		status = tester.AssertionResultStatusSuccess
	}

	return tester.NewAssertionResultList(
		tester.AssertionResult{
			Name:     a.name,
			Status:   status,
			Expected: a.asserter.GetExpected(),
			Actual:   a.asserter.GetActual(),
		},
	)
}
