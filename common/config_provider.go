package common

import (
	"go.uber.org/config"
	"os"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=mocks/$GOFILE

type ConfigProviderFactory interface {
	CreateByFiles(filenames ...string) (ConfigProvider, error)
	CreateByOptions(options ...config.YAMLOption) (ConfigProvider, error)
}

func NewFxConfigProviderFactory() ConfigProviderFactory {
	return new(configProviderFactory)
}

type configProviderFactory struct{}

func (f *configProviderFactory) CreateByOptions(options ...config.YAMLOption) (ConfigProvider, error) {
	provider, err := config.NewYAML(append(
		[]config.YAMLOption{
			config.Expand(os.LookupEnv),
		},
		options...,
	)...)
	if err != nil {
		return nil, err
	}

	return &configProvider{provider}, nil
}

func (f *configProviderFactory) CreateByFiles(filenames ...string) (ConfigProvider, error) {
	opts := make([]config.YAMLOption, len(filenames))
	for i, filename := range filenames {
		opts[i] = config.File(filename)
	}

	return f.CreateByOptions(opts...)
}

type ConfigProvider interface {
	Populate(interface{}) error
	PopulateByKey(string, interface{}) error
}

type configProvider struct {
	provider config.Provider
}

func (p *configProvider) Populate(target interface{}) error {
	return p.PopulateByKey("", target)
}

func (p *configProvider) PopulateByKey(key string, target interface{}) error {
	return p.provider.Get(key).Populate(target)
}
