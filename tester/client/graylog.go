package client

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester/client/graylog"
	"github.com/byorty/contractor/tester/client/graylog/saved"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
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
