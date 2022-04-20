package e2e

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester/client"
)

type TestRunContainer map[string]*TestRun

type TestRun struct {
	Messages common.Map[string, common.List[client.GraylogMessage]]
}
