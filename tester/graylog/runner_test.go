package graylog_test

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/graylog"
	"github.com/byorty/contractor/tester/graylog/client/models"
	"github.com/byorty/contractor/tester/graylog/client/saved"
	gm "github.com/byorty/contractor/tester/graylog/mocks"
	tm "github.com/byorty/contractor/tester/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGraylogSuite(t *testing.T) {
	suite.Run(t, new(GraylogSuite))
}

type GraylogSuite struct {
	common.Suite
	graylogClient  *gm.MockClient
	runner         tester.Runner
	correlationId1 string
	correlationId2 string
	graylogResp1   *saved.SearchRelativeOK
	graylogResp2   *saved.SearchRelativeOK
	graylogErr     error
	assertResult1  tester.AssertionResult
	assertResult2  tester.AssertionResult
	testCase       tester.TestCase2
}

func (s *GraylogSuite) SetupTest() {
	s.Suite.SetupTest()
	s.testCase = tester.TestCase2{
		Name: randomdata.Alphanumeric(32),
	}
	s.graylogClient = gm.NewMockClient(s.Ctrl)
	s.runner = graylog.NewRunner(
		context.Background(),
		common.NewFxLoggerFactory().CreateCommonLogger(),
		s.graylogClient,
		nil,
	)
	s.runner.Setup(s.testCase)

	s.correlationId1 = randomdata.Alphanumeric(16)
	s.correlationId2 = randomdata.Alphanumeric(16)
	s.graylogErr = errors.Errorf("network error")
	s.graylogResp1 = &saved.SearchRelativeOK{
		Payload: &models.SearchResponse{
			Messages: []*models.SearchResponseMessagesItems0{
				{
					Message: map[string]interface{}{
						"correlation_id": s.correlationId1,
						"msg":            randomdata.City(),
					},
				},
			},
		},
	}
	s.graylogResp2 = &saved.SearchRelativeOK{
		Payload: &models.SearchResponse{
			Messages: []*models.SearchResponseMessagesItems0{
				{
					Message: map[string]interface{}{
						"correlation_id": s.correlationId1,
						"msg":            "start",
					},
				},
				{
					Message: map[string]interface{}{
						"correlation_id": s.correlationId1,
						"msg":            randomdata.City(),
					},
				},
				{
					Message: map[string]interface{}{
						"correlation_id": s.correlationId2,
						"msg":            "start",
					},
				},
				{
					Message: map[string]interface{}{
						"correlation_id": s.correlationId2,
						"msg":            randomdata.City(),
					},
				},
				{
					Message: map[string]interface{}{
						"correlation_id": s.correlationId2,
						"msg":            randomdata.Email(),
					},
				},
			},
		},
	}
	s.assertResult1 = tester.AssertionResult{
		Status:   tester.AssertionResultStatusFailure,
		Expected: randomdata.Alphanumeric(8),
		Actual:   randomdata.Alphanumeric(8),
	}
	s.assertResult2 = tester.AssertionResult{
		Status:   tester.AssertionResultStatusSuccess,
		Expected: randomdata.Alphanumeric(8),
		Actual:   randomdata.Alphanumeric(8),
	}
}

func (s *GraylogSuite) TestRun() {
	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(nil, s.graylogErr).
		Times(1)
	list := s.runner.Run(tester.Asserter2List{})
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal(s.graylogErr.Error(), list.Get(0).Assertions.Get(0).Actual)

	assertion1 := tm.NewMockAssertion2(s.Ctrl)
	assertion2 := tm.NewMockAssertion2(s.Ctrl)
	assertions := tester.NewAsserter2List(assertion1, assertion2)
	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(&saved.SearchRelativeOK{
			Payload: &models.SearchResponse{
				Messages: []*models.SearchResponseMessagesItems0{},
			},
		}, nil).
		Times(1)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal("runner messages present", list.Get(0).Assertions.Get(0).Expected)
	s.Equal("nil", list.Get(0).Assertions.Get(0).Actual)

	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(s.graylogResp1, nil).
		Times(1)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal(fmt.Sprintf("%s#%s", s.testCase.Name, s.correlationId1), list.Get(0).Name)
	s.Equal("runner messages present", list.Get(0).Assertions.Get(0).Expected)
	s.Equal("nil", list.Get(0).Assertions.Get(0).Actual)

	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(s.graylogResp2, nil).Times(2)
	assertion1.EXPECT().Assert(gomock.Any()).Return(tester.NewAssertionResultList(s.assertResult1)).Times(2)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 2)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal(fmt.Sprintf("%s#%s", s.testCase.Name, s.correlationId1), list.Get(0).Name)
	s.Equal(s.assertResult1.Expected, list.Get(0).Assertions.Get(0).Expected)
	s.Len(list.Get(1).Assertions.Entries(), 1)
	s.Equal(fmt.Sprintf("%s#%s", s.testCase.Name, s.correlationId2), list.Get(1).Name)
	s.Equal(s.assertResult1.Expected, list.Get(1).Assertions.Get(0).Expected)

	assertion1.EXPECT().Assert(s.graylogResp2.Payload.Messages[1].Message.(map[string]interface{})["msg"]).Return(tester.NewAssertionResultList(s.assertResult1)).Times(1)
	assertion1.EXPECT().Assert(s.graylogResp2.Payload.Messages[3].Message.(map[string]interface{})["msg"]).Return(tester.NewAssertionResultList(s.assertResult2)).Times(1)
	assertion2.EXPECT().Assert(gomock.Any()).Return(tester.NewAssertionResultList(s.assertResult2)).Times(1)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 2)
	s.Equal(s.assertResult2.Expected, list.Get(0).Assertions.Get(0).Expected)
	s.Equal(s.assertResult2.Expected, list.Get(0).Assertions.Get(1).Expected)
}
