package tester

import (
	"github.com/byorty/contractor/common"
	"time"
)

type TestCase2List struct {
	common.List[TestCase2]
}

func NewTestCase2List() TestCase2List {
	return TestCase2List{common.NewList[TestCase2]()}
}

func (l *TestCase2List) ImportFromMap(testCaseMap common.Map[string, TestCase2]) {
	for name, testCase := range testCaseMap.Entries() {
		testCase.Name = name
		l.Add(testCase)
	}
}

type TestCase2 struct {
	Name       string         `yaml:"name"`
	Tags       []string       `yaml:"tags"`
	Priority   int            `yaml:"priority"`
	Setup      TestCase2Setup `yaml:"setup"`
	Assertions []Assertion2   `yaml:"assertions"`
}

func (t *TestCase2) ContainsTags(expectedTags []string) bool {
	if len(expectedTags) == 0 {
		return true
	}

	for _, expectedTag := range expectedTags {
		for _, actualTag := range t.Tags {
			if expectedTag == actualTag {
				return true
			}
		}
	}

	return false
}

type TestCase2Setup struct {
	Query        string                 `yaml:"query"`
	Range        time.Duration          `yaml:"range"`
	Trigger      string                 `yaml:"trigger"`
	Headers      map[string]interface{} `yaml:"headers"`
	Parameters   map[string]interface{} `yaml:"parameters"`
	HeaderParams map[string]interface{} `yaml:"-"`
	PathParams   map[string]interface{} `yaml:"-"`
	QueryParams  map[string]interface{} `yaml:"-"`
	CookieParams map[string]interface{} `yaml:"-"`
}

type Assertion2 struct {
	Name   string      `yaml:"name"`
	Type   string      `yaml:"type"`
	Assert interface{} `yaml:"assert"`
}
