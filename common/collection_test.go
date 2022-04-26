package common_test

import (
	"github.com/byorty/contractor/common"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCollectionSuite(t *testing.T) {
	suite.Run(t, new(CollectionSuite))
}

type CollectionSuite struct {
	suite.Suite
}

func (s *CollectionSuite) TestListSorting() {
	list := common.NewListFromSlice[int](10, 2, 14, 1, -1, 22, 9)
	list.Sort(func(a, b int) bool {
		return a < b
	})
	s.Equal(-1, list.Get(0))
	s.Equal(1, list.Get(1))
	s.Equal(2, list.Get(2))
	s.Equal(9, list.Get(3))
	s.Equal(10, list.Get(4))
	s.Equal(14, list.Get(5))
	s.Equal(22, list.Get(6))
}
