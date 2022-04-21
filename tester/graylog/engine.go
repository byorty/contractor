package graylog

import (
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	tc "github.com/byorty/contractor/tester/client"
	"github.com/go-openapi/runtime/client"
	"go.uber.org/config"
	"io/ioutil"
)

func NewFxEngine(
	logger common.Logger,
	configProviderFactory common.ConfigProviderFactory,
	graylogClient tc.GraylogClient,
) tester.Engine {
	return &engine{
		logger:                logger.Named("graylog_engine"),
		configProviderFactory: configProviderFactory,
		graylogClient:         graylogClient,
		testCases:             tester.NewTestCase2List(),
		cfg:                   new(Config),
	}
}

type engine struct {
	logger                common.Logger
	configProviderFactory common.ConfigProviderFactory
	graylogClient         tc.GraylogClient
	testCases             tester.TestCase2List
	cfg                   *Config
}

func (e *engine) Configure(data interface{}) error {
	configProvider, err := e.configProviderFactory.CreateByOptions(config.Static(data))
	if err != nil {
		e.logger.Error(err)
		return err
	}

	err = configProvider.Populate(e.cfg)
	if err != nil {
		e.logger.Error(err)
		return err
	}

	files, err := ioutil.ReadDir(e.cfg.SpecDir)
	if err != nil {
		e.logger.Error(err)
		return err
	}

	filenames := common.NewList[string]()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filenames.Add(file.Name())
	}

	testCaseProvider, err := e.configProviderFactory.CreateByFiles(filenames.Entries()...)
	if err != nil {
		e.logger.Error(err)
		return err
	}

	testCaseMap := common.NewMap[string, tester.TestCase2]()
	err = testCaseProvider.Populate(testCaseMap)
	if err != nil {
		e.logger.Error(err)
		return err
	}

	e.testCases.ImportFromMap(testCaseMap)
	return nil
}

func (e *engine) GetTestCase2List() tester.TestCase2List {
	return e.testCases
}

func (e *engine) CreateRunner() tester.Runner {
	return NewRunner(
		e.logger,
		e.graylogClient,
		client.BasicAuth(e.cfg.Username, e.cfg.Password),
	)
}
