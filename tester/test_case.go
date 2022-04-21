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
	Name string
	Tags []string
	//PostProcessors []common.PostProcessor
	Priority   int
	Setup      TestCase2Setup
	Assertions []Assertion2
}

type TestCase2Setup struct {
	Query string
	Range time.Duration
}

type Assertion2 struct {
	Name   string
	Type   string
	Assert interface{}
}
