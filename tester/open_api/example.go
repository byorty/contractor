package open_api

import "github.com/byorty/contractor/common"

const (
	xExamples       = "x-examples"
	xTagsName       = "x-tags"
	xPriority       = "x-priority"
	xPostProcessors = "x-post-processors"
)

type XOperationRequestExample struct {
	Headers    map[string]interface{} `json:"headers"`
	Parameters map[string]interface{} `json:"parameters"`
}

type XOperationResponseExample struct {
	StatusCode int                    `json:"status_code"`
	Body       map[string]interface{} `json:"body"`
}

type XOperationExample struct {
	Request        XOperationRequestExample  `json:"request"`
	Response       XOperationResponseExample `json:"response"`
	Tags           []string                  `json:"tags"`
	Priority       int                       `json:"priority"`
	PostProcessors []common.PostProcessor    `json:"post_processors"`
}

type XOperationExamples map[string]XOperationExample
