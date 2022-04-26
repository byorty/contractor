package graylog

import (
	"github.com/byorty/contractor/tester/graylog/client"
	"github.com/byorty/contractor/tester/graylog/client/saved"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"net/url"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type Client interface {
	SearchRelative(params *saved.SearchRelativeParams, authInfo runtime.ClientAuthInfoWriter, opts ...saved.ClientOption) (*saved.SearchRelativeOK, error)
}

type graylogClient struct {
	saved.ClientService
}

func NewClient(
	cfg *Config,
) (Client, error) {
	u, err := url.Parse(cfg.Url)
	if err != nil {
		return nil, err
	}

	generateClient := client.NewHTTPClientWithConfig(strfmt.Default, &client.TransportConfig{
		Host:     u.Host,
		BasePath: u.Path,
		Schemes:  []string{u.Scheme},
	})

	return &graylogClient{
		ClientService: generateClient.Saved,
	}, nil
}
