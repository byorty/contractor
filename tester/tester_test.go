package tester_test

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	"github.com/stretchr/testify/suite"
	"sort"
	"testing"
)

func TestTesterTestSuite(t *testing.T) {
	suite.Run(t, new(TesterTestSuite))
}

type TesterTestSuite struct {
	suite.Suite
	container tester.TestCaseContainer
}

func (s *TesterTestSuite) TestContainerSorting() {
	container := tester.TestCaseContainer{
		{
			Name:     "ABC",
			Template: &common.Template{},
		},
		{
			Name:     "ABS",
			Template: &common.Template{},
		},
		{
			Name:     "LMN",
			Template: &common.Template{},
		},
		{
			Name: "BCD",
			Template: &common.Template{
				Priority: 1,
			},
		},
		{
			Name: "BCD2",
			Template: &common.Template{
				Priority: 1,
			},
		},
	}

	sort.Sort(container)

	s.Equal("BCD", container[0].Name)
	s.Equal("BCD2", container[1].Name)
	s.Equal("ABC", container[2].Name)
	s.Equal("ABS", container[3].Name)
	s.Equal("LMN", container[4].Name)
}
