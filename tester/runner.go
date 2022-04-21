package tester

import "github.com/byorty/contractor/common"

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type Runner interface {
	Setup(testCase TestCase2)
	Run(assertions Asserter2List) RunnerReportList
}

type RunnerReportList struct {
	common.List[RunnerReport]
}

func NewRunnerReportList() RunnerReportList {
	return RunnerReportList{common.NewList[RunnerReport]()}
}

type RunnerReportDetail struct {
	Name string
	Data interface{}
}

type RunnerReport struct {
	Name       string
	Details    []RunnerReportDetail
	Assertions AssertionResultList
}
