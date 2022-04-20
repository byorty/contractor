package converter_test

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestOA2ConverterTestSuite(t *testing.T) {
	suite.Run(t, new(OA2ConverterTestSuite))
}

type OA2ConverterTestSuite struct {
	ConverterTestSuite
}

func (s *OA2ConverterTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.converter = converter.NewFxOa2Converter()
	s.arguments = common.Arguments{
		SpecLocation: s.makeSpecFilename("oa2.yml"),
	}
}

func (s *OA2ConverterTestSuite) Test() {
	s.T().Log(s.arguments.SpecLocation)
}
