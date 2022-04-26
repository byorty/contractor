package graylog

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester"
	rc "github.com/go-openapi/runtime/client"
	"go.uber.org/config"
	"io/ioutil"
	"path/filepath"
	"time"
)

func NewFxEngine(
	ctx context.Context,
	logger common.Logger,
	configProviderFactory common.ConfigProviderFactory,
) tester.Engine {
	return &engine{
		ctx:                   ctx,
		logger:                logger.Named("graylog_engine"),
		configProviderFactory: configProviderFactory,
		testCases:             tester.NewTestCase2List(),
		cfg:                   new(Config),
	}
}

type engine struct {
	ctx                   context.Context
	logger                common.Logger
	configProviderFactory common.ConfigProviderFactory
	graylogClient         Client
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

	e.graylogClient, err = NewClient(e.cfg)
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

		filenames.Add(filepath.Join(e.cfg.SpecDir, file.Name()))
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
	ctx, _ := context.WithDeadline(e.ctx, time.Now().Add(e.cfg.Timeout))
	return NewRunner(
		ctx,
		e.logger,
		e.graylogClient,
		rc.BasicAuth(e.cfg.Username, e.cfg.Password),
	)
}
