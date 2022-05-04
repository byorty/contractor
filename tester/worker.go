package tester

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/tester/open_api/reader"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var (
	ErrTestHasFailureStatus = errors.New("Test Is Failed")
)

type WorkerIn struct {
	fx.In
	ConverterContainer reader.Container
	Tester             Tester
	Reporters          []Reporter `group:"reporter"`
}

func NewFxWorker(in WorkerIn) common.Worker {
	return &worker{
		converterContainer: in.ConverterContainer,
		tester:             in.Tester,
		reporters:          in.Reporters,
	}
}

type worker struct {
	converterContainer reader.Container
	tester             Tester
	reporters          []Reporter
}

func (w *worker) GetType() common.WorkerType {
	return common.WorkerTypeTest
}

func (w *worker) Configure(ctx context.Context, arguments common.Arguments) error {
	crt, err := w.converterContainer.Get(arguments.SpecType)
	if err != nil {
		return err
	}

	templateContainers, err := crt.Convert(ctx, arguments)
	if err != nil {
		return err
	}

	w.tester.Configure(ctx, arguments, templateContainers)

	return nil
}

func (w *worker) Run() error {
	container, err := w.tester.Test()
	if err != nil {
		return err
	}

	for _, reporter := range w.reporters {
		err := reporter.Report(container)
		if err != nil {
			return err
		}
	}

	if container.HasError() {
		return ErrTestHasFailureStatus
	}

	return nil
}
