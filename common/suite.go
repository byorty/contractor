package common

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	Ctrl *gomock.Controller
}

func (s *Suite) SetupTest() {
	s.Ctrl = gomock.NewController(s.T())
}

func (s *Suite) TearDownTest() {
	s.Ctrl.Finish()
}
