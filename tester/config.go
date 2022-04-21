package tester

import "github.com/byorty/contractor/common"

func NewFxConfig(
	arguments common.Arguments,
	logger common.Logger,
	configProviderFactory common.ConfigProviderFactory,
) (*Config, error) {
	configProvider, err := configProviderFactory.CreateByFiles(arguments.SpecLocation)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	cfg := new(Config)
	err = configProvider.Populate(cfg)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return cfg, nil
}

type Config struct {
	Testers   map[string]interface{}
	Reporters map[string]interface{}
	Variables map[string]interface{}
}
