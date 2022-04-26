package common_test

import (
	"github.com/byorty/contractor/common"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestDataCrawlerSuite(t *testing.T) {
	suite.Run(t, new(DataCrawlerTestSuite))
}

type DataCrawlerTestSuite struct {
	suite.Suite
}

func (s *DataCrawlerTestSuite) TestWalk() {
	data := map[string]interface{}{
		"$.k1":  "a",
		"k2":    1,
		"slice": []interface{}{1, 2, 3},
		"sliceMap": []interface{}{
			map[string]interface{}{
				"k1": "b",
				"k2": 2,
			},
			map[string]interface{}{
				"k1": "c",
				"k2": 3,
			},
		},
		"map": map[string]interface{}{
			"k1":    "d",
			"k2":    4,
			"slice": []interface{}{1, 2},
			"slice2": []interface{}{
				map[string]interface{}{
					"k1": 10,
				},
			},
		},
	}
	flatData := make(map[string]interface{})
	crawler := common.NewFxDataCrawler()
	crawler.Walk(data, func(k string, v interface{}) {
		flatData[k] = v
	}, common.WithPrefix("$"), common.WithJoinKeys(), common.WithSkipCollections())

	s.Len(flatData, 14)
	s.Equal(data["$.k1"], flatData["$.k1"])
	s.Equal(data["k2"], flatData["$.k2"])
	s.Equal(data["slice"].([]interface{})[0], flatData["$.slice[0]"])
	s.Equal(data["slice"].([]interface{})[1], flatData["$.slice[1]"])
	s.Equal(data["slice"].([]interface{})[2], flatData["$.slice[2]"])
	s.Equal(data["sliceMap"].([]interface{})[0].(map[string]interface{})["k1"], flatData["$.sliceMap[0].k1"])
	s.Equal(data["sliceMap"].([]interface{})[0].(map[string]interface{})["k2"], flatData["$.sliceMap[0].k2"])
	s.Equal(data["sliceMap"].([]interface{})[1].(map[string]interface{})["k1"], flatData["$.sliceMap[1].k1"])
	s.Equal(data["sliceMap"].([]interface{})[1].(map[string]interface{})["k2"], flatData["$.sliceMap[1].k2"])
	s.Equal(data["map"].(map[string]interface{})["k1"], flatData["$.map.k1"])
	s.Equal(data["map"].(map[string]interface{})["k2"], flatData["$.map.k2"])
	s.Equal(data["map"].(map[string]interface{})["slice"].([]interface{})[0], flatData["$.map.slice[0]"])
	s.Equal(data["map"].(map[string]interface{})["slice"].([]interface{})[1], flatData["$.map.slice[1]"])
	s.Equal(data["map"].(map[string]interface{})["slice2"].([]interface{})[0].(map[string]interface{})["k1"], flatData["$.map.slice2[0].k1"])
}
