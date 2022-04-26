package common

func NewFxConfig(
	arguments Arguments,
	loggerFactory LoggerFactory,
	configProviderFactory ConfigProviderFactory,
) (*Config, error) {
	logger := loggerFactory.CreateCommonLogger().Named("config_tester")
	configProvider, err := configProviderFactory.CreateByFiles(arguments.ConfigFilename)
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
	Testers   map[string]interface{} `yaml:"testers"`
	Reporters map[string]interface{} `yaml:"reporters"`
	Variables map[string]interface{} `yaml:"variables"`
}
