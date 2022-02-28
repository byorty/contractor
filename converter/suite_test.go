package converter_test

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
)

type ConverterTestSuite struct {
	suite.Suite
	ctx       context.Context
	converter converter.Converter
	arguments common.Arguments
}

func (s ConverterTestSuite) makeSpecFilename(filename string) string {
	dir, err := os.Getwd()
	s.Nil(err)
	return filepath.Join(dir, "..", "specs", filename)
}
