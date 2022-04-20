package client

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester/client/graylog"
	"github.com/byorty/contractor/tester/client/graylog/saved"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"time"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type GraylogClient interface {
	SearchRelative(params *saved.SearchRelativeParams, authInfo runtime.ClientAuthInfoWriter, opts ...saved.ClientOption) (*saved.SearchRelativeOK, error)
}

type graylogClient struct {
	saved.ClientService
}

func NewFxGraylogClient(
	args common.Arguments,
) GraylogClient {
	generateClient := graylog.NewHTTPClientWithConfig(strfmt.Default, &graylog.TransportConfig{
		Host:     args.Url,
		BasePath: "api",
		Schemes:  []string{"https"},
	})

	return &graylogClient{
		ClientService: generateClient.Saved,
	}
}

type GraylogMessage struct {
	App           string `json:"app"`
	CorrelationId string `json:"correlation_id"`
	Target        string `json:"target"`
	DateTime      time.Time
	Message       string `json:"message"`
	Method        string `json:"method"`
}

type GraylogMessages []GraylogMessage

func (m GraylogMessages) Len() int {
	return len(m)
}

func (m GraylogMessages) Less(i, j int) bool {
	return m[i].DateTime.UnixNano() > m[j].DateTime.UnixNano()
}

func (m GraylogMessages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
