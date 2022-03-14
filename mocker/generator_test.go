package mocker_test

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/mocker"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}

type GeneratorTestSuite struct {
	suite.Suite
}

func (g *GeneratorTestSuite) TestEqGenerator() {
	expected := randomdata.Alphanumeric(16)
	gen := mocker.NewEqGenerator(expected)
	g.Equal(expected, gen.Generate())
}

func (g *GeneratorTestSuite) TestPositiveGenerator() {
	gen := mocker.NewPositiveGenerator()
	g.Less(0, gen.Generate())
}

func (g *GeneratorTestSuite) TestMinGenerator() {
	gen := mocker.NewMinGenerator(100)
	g.Less(0, gen.Generate())
}

func (g *GeneratorTestSuite) TestMaxGenerator() {
	gen := mocker.NewMaxGenerator(100)
	g.Greater(101, gen.Generate())
}

func (g *GeneratorTestSuite) TestRangeGenerator() {
	gen := mocker.NewRangeGenerator(1, 99)
	actual := gen.Generate()
	g.Less(0, actual)
	g.Greater(100, actual)
}

func (g *GeneratorTestSuite) TestRegexGenerator() {
	exp := "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	gen := mocker.NewRegexGenerator(exp)

	matched, err := regexp.MatchString(exp, gen.Generate().(string))
	g.Nil(err)
	g.True(matched)
}

func (g *GeneratorTestSuite) TestEmptyGenerator() {
	gen := mocker.NewEmptyGenerator()
	g.Empty(gen.Generate())
}

func (g *GeneratorTestSuite) TestDateGenerator() {
	gen := mocker.NewDateGenerator(common.ExpressionDateLayoutRFC3339)
	_, err := time.Parse(time.RFC3339, gen.Generate().(string))
	g.Nil(err)

	gen = mocker.NewDateGenerator(common.ExpressionDateLayoutRFC3339Nano)
	_, err = time.Parse(time.RFC3339Nano, gen.Generate().(string))
	g.Nil(err)

	_, err = time.Parse(time.Layout, gen.Generate().(string))
	g.NotNil(err)
}

func (g *GeneratorTestSuite) TestContainsGenerator() {
	sub := randomdata.Alphanumeric(32)
	gen := mocker.NewContainsGenerator(sub)
	g.Contains(gen.Generate(), sub)
}
