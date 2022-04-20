package tester

import "github.com/byorty/contractor/common"

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type TestRunner interface {
	Setup(name string, testCase TestCaseDefinition)
	Run(assertions Assertion2List) TestRunnerReportList
}

type TestRunnerReportList struct {
	common.List[TestRunnerReport]
}

func NewTestRunnerReportList() TestRunnerReportList {
	return TestRunnerReportList{common.NewList[TestRunnerReport]()}
}

type TestRunnerReport struct {
	Name       string
	Assertions AssertionResultList
}
