package tester

import (
	"github.com/byorty/contractor/common"
	"time"
)

type TestCaseDefinitionMap common.Map[string, TestCaseDefinition]

type TestCaseDefinition struct {
	Tags []string
	//PostProcessors []common.PostProcessor
	Priority   int
	Setup      SetupTestCaseDefinition
	Assertions []AssertionDefinition
}

type SetupTestCaseDefinition struct {
	Query string
	Range time.Duration
}

type AssertionDefinition struct {
	Name   string
	Type   string
	Assert interface{}
}
