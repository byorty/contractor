package tester

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
)

func NewFxWorker(
	cnv converter.Converter,
	tester Tester,
) common.Worker {
	return &Worker{
		converter: cnv,
		tester:    tester,
	}
}

type Worker struct {
	converter converter.Converter
	tester    Tester
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
	err := w.tester.Test()
	if err != nil {
		return err
	}

	return nil
}
