package runner_test

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/client/graylog/models"
	"github.com/byorty/contractor/tester/client/graylog/saved"
	cm "github.com/byorty/contractor/tester/client/mocks"
	tm "github.com/byorty/contractor/tester/mocks"
	"github.com/byorty/contractor/tester/runner"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGraylogSuite(t *testing.T) {
	suite.Run(t, new(GraylogSuite))
}

type GraylogSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	assertionFactory *tm.MockAssertionFactory
	graylogClient    *cm.MockGraylogClient
	runner           tester.TestRunner
	runnerName       string
	correlationId1   string
	correlationId2   string
	graylogResp1     *saved.SearchRelativeOK
	graylogResp2     *saved.SearchRelativeOK
	assertResult1    tester.AssertionResult
	assertResult2    tester.AssertionResult
}

func (s *GraylogSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.assertionFactory = tm.NewMockAssertionFactory(s.ctrl)
	s.assertionFactory = tm.NewMockAssertionFactory(s.ctrl)
	s.graylogClient = cm.NewMockGraylogClient(s.ctrl)
	s.runnerName = randomdata.Alphanumeric(32)
	s.runner = runner.NewFxGraylog(
		common.NewFxLoggerFactory().CreateCommonLogger(),
		s.assertionFactory,
		s.graylogClient,
	)
	s.runner.Setup(s.runnerName, tester.TestCaseDefinition{})

	s.correlationId1 = randomdata.Alphanumeric(16)
	s.correlationId2 = randomdata.Alphanumeric(16)
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

func (s *GraylogSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *GraylogSuite) TestRun() {
	err := errors.Errorf("network error")
	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(nil, err).
		Times(1)
	list := s.runner.Run(tester.Assertion2List{})
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal(err.Error(), list.Get(0).Assertions.Get(0).Actual)

	assertion1 := tm.NewMockAssertion2(s.ctrl)
	assertion2 := tm.NewMockAssertion2(s.ctrl)
	assertions := tester.NewAssertion2List(assertion1, assertion2)
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
	s.Equal("graylog messages present", list.Get(0).Assertions.Get(0).Expected)
	s.Equal("nil", list.Get(0).Assertions.Get(0).Actual)

	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(s.graylogResp1, nil).
		Times(1)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal(fmt.Sprintf("%s#%s", s.runnerName, s.correlationId1), list.Get(0).Name)
	s.Equal("graylog messages present", list.Get(0).Assertions.Get(0).Expected)
	s.Equal("nil", list.Get(0).Assertions.Get(0).Actual)

	s.graylogClient.
		EXPECT().
		SearchRelative(gomock.Any(), gomock.Any()).
		Return(s.graylogResp2, nil).Times(2)
	assertion1.EXPECT().Assert(gomock.Any()).Return(tester.NewAssertionResultList(s.assertResult1), nil).Times(2)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 2)
	s.Len(list.Get(0).Assertions.Entries(), 1)
	s.Equal(fmt.Sprintf("%s#%s", s.runnerName, s.correlationId1), list.Get(0).Name)
	s.Equal(s.assertResult1.Expected, list.Get(0).Assertions.Get(0).Expected)
	s.Len(list.Get(1).Assertions.Entries(), 1)
	s.Equal(fmt.Sprintf("%s#%s", s.runnerName, s.correlationId2), list.Get(1).Name)
	s.Equal(s.assertResult1.Expected, list.Get(1).Assertions.Get(0).Expected)

	assertion1.EXPECT().Assert(s.graylogResp2.Payload.Messages[1].Message.(map[string]interface{})["msg"]).Return(tester.NewAssertionResultList(s.assertResult1), nil).Times(1)
	assertion1.EXPECT().Assert(s.graylogResp2.Payload.Messages[3].Message.(map[string]interface{})["msg"]).Return(tester.NewAssertionResultList(s.assertResult2), nil).Times(1)
	assertion2.EXPECT().Assert(gomock.Any()).Return(tester.NewAssertionResultList(s.assertResult2), nil).Times(1)
	list = s.runner.Run(assertions)
	s.Len(list.Entries(), 1)
	s.Len(list.Get(0).Assertions.Entries(), 2)
	s.Equal(s.assertResult2.Expected, list.Get(0).Assertions.Get(0).Expected)
	s.Equal(s.assertResult2.Expected, list.Get(0).Assertions.Get(1).Expected)
}
