package tester

import (
	"context"
	"github.com/byorty/contractor/common"
)

func NewFxWorker2(
	cfg *Config,
	logger common.Logger,
	configProviderFactory common.ConfigProviderFactory,
	testEngineFactory EngineFactory,
	assertionFactory Assertion2Factory,
	reporterFactory Reporter2Factory,
) common.Worker {
	return &worker2{
		cfg:                   cfg,
		logger:                logger.Named("worker"),
		configProviderFactory: configProviderFactory,
		testEngineFactory:     testEngineFactory,
		assertionFactory:      assertionFactory,
		reporterFactory:       reporterFactory,
		engines:               common.NewList[Engine](),
		reporters:             common.NewList[Reporter2](),
	}
}

type worker2 struct {
	cfg                   *Config
	logger                common.Logger
	configProviderFactory common.ConfigProviderFactory
	testEngineFactory     EngineFactory
	assertionFactory      Assertion2Factory
	reporterFactory       Reporter2Factory
	engines               common.List[Engine]
	reporters             common.List[Reporter2]
}

func (w *worker2) GetType() common.WorkerType {
	return common.WorkerTypeTest2
}

func (w *worker2) Configure(ctx context.Context, arguments common.Arguments) error {
	for engineName, engineConfig := range w.cfg.Testers {
		engine, err := w.testEngineFactory.Create(engineName)
		if err != nil {
			w.logger.Error(err)
			return err
		}

		err = engine.Configure(engineConfig)
		if err != nil {
			w.logger.Error(err)
			return err
		}

		w.engines.Add(engine)
	}

	for reporterName, reporterConfig := range w.cfg.Reporters {
		reporter, err := w.reporterFactory.Create(reporterName, reporterConfig)
		if err != nil {
			w.logger.Error(err)
			return err
		}

		w.reporters.Add(reporter)
	}

	return nil
}

func (w *worker2) Run() error {
	testRunReports := NewRunnerReportList()
	for _, engine := range w.engines.Entries() {
		for _, testCase := range engine.GetTestCase2List().Entries() {
			asserters := NewAsserter2List()
			for _, assertion := range testCase.Assertions {
				asserter, err := w.assertionFactory.Create(assertion.Name, assertion.Assert)
				if err != nil {
					w.logger.Error(err)
					return err
				}

				asserters.Add(asserter)
			}

			testCaseRunner := engine.CreateRunner()
			testCaseRunner.Setup(testCase)
			reports := testCaseRunner.Run(asserters)
			testRunReports.Add(reports.Entries()...)
		}
	}

	for _, reporter := range w.reporters.Entries() {
		for _, report := range testRunReports.Entries() {
			reporter.Report(report)
		}
	}

	return nil
}
