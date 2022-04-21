package assertion_test

import (
	"errors"
	"github.com/byorty/contractor/common"
	cm "github.com/byorty/contractor/common/mocks"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/assertion"
	tm "github.com/byorty/contractor/tester/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestJsonContainsSuite(t *testing.T) {
	suite.Run(t, new(JsonContainsSuite))
}

type JsonContainsSuite struct {
	common.Suite
	assertion            tester.Asserter2
	expressionFactory    *cm.MockExpressionFactory
	expressionFactoryErr error
	asserter             *tm.MockAsserter
	asserterErr          error
	definition           map[string]interface{}
	jsonObject           string
}

func (s *JsonContainsSuite) SetupTest() {
	s.Suite.SetupTest()
	s.asserter = tm.NewMockAsserter(s.Ctrl)
	s.asserterErr = errors.New("asserter error")
	s.expressionFactory = cm.NewMockExpressionFactory(s.Ctrl)
	s.expressionFactoryErr = errors.New("expression factory error")
	s.definition = map[string]interface{}{
		"status":                     "eq('SUBSCRIPTION_PURCHASE_STATUS_ACTIVE')",
		"status_2":                   "eq('SUBSCRIPTION_PURCHASE_STATUS_ACTIVE_2')",
		"auto_payment":               "eq('SUBSCRIPTION_PURCHASE_AUTO_PAYMENT_ACTIVE')",
		"expiration_at":              "date_range('RFC3339', -1000h, -1min)",
		"auto_payment_attempt_count": "max(2)",
	}
	s.assertion = assertion.NewJsonContains(
		common.NewFxLoggerFactory().CreateCommonLogger(),
		common.NewFxDataCrawler(),
		s.expressionFactory,
		s.definition,
	)
	s.jsonObject = `{
	"status": "SUBSCRIPTION_PURCHASE_STATUS_ACTIVE",
	"auto_payment": "SUBSCRIPTION_PURCHASE_AUTO_PAYMENT_ACTIVE",
	"expiration_at": "1970-01-01T00:00:00Z",
	"auto_payment_attempt_count": 1
}`
}

func (s *JsonContainsSuite) TestAssert() {
	results := s.assertion.Assert(nil)
	s.Len(results.Entries(), 1)
	s.Equal(tester.AssertionResultStatusFailure, results.Get(0).Status)
	s.Equal("json present", results.Get(0).Expected)

	s.expressionFactory.EXPECT().Create(common.ExpressionTypeAsserter, s.definition["status"]).Return(nil, s.expressionFactoryErr).Times(1)
	s.expressionFactory.EXPECT().Create(common.ExpressionTypeAsserter, gomock.Any()).Return(s.asserter, nil).Times(4)
	s.asserter.EXPECT().Assert("SUBSCRIPTION_PURCHASE_AUTO_PAYMENT_ACTIVE").Return(s.asserterErr).Times(1)
	s.asserter.EXPECT().Assert(gomock.Any()).Return(nil).Times(2)
	s.asserter.EXPECT().GetExpected().Return("expected value").Times(4)
	s.asserter.EXPECT().GetActual().Return("actual value").Times(3)
	results = s.assertion.Assert(s.jsonObject)
	s.Len(results.Entries(), 5)
	var successCount, failureCount int
	for i := 0; i < results.Len(); i++ {
		if results.Get(i).Status == tester.AssertionResultStatusSuccess {
			successCount++
		} else {
			failureCount++
		}
	}
	s.Equal(2, successCount)
	s.Equal(3, failureCount)
}
