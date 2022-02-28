package common_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestMediaConverterTestSuite(t *testing.T) {
	suite.Run(t, new(MediaConverterTestSuite))
}

type MediaConverterTestSuite struct {
	suite.Suite
}
