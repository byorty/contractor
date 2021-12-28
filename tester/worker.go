package tester

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"go.uber.org/fx"
)

type WorkerIn struct {
	fx.In
	Converter converter.Converter
	Tester    Tester
	Reporters []Reporter `group:"reporter"`
}

func NewFxWorker(in WorkerIn) common.Worker {
	return &Worker{
		converter: in.Converter,
		tester:    in.Tester,
		reporters: in.Reporters,
	}
}

type Worker struct {
	converter converter.Converter
	tester    Tester
	reporters []Reporter
}

func (w *Worker) GetType() common.WorkerKind {
	return common.WorkerKindTest
}

func (w *Worker) Configure(ctx context.Context, arguments common.Arguments) error {
	templateContainers, err := w.converter.Convert(ctx, arguments)
	if err != nil {
		return err
	}

	w.tester.Configure(ctx, templateContainers)

	return nil
}

func (w *Worker) Run() error {
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

	return nil
}
