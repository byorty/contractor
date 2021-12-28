package mocker

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"go.uber.org/fx"
)

type WorkerIn struct {
	fx.In
	Ctx       context.Context
	Converter converter.Converter
	Mocker    Mocker
}

func NewFxWorker(in WorkerIn) common.Worker {
	return &worker{
		ctx:       in.Ctx,
		converter: in.Converter,
		mocker:    in.Mocker,
	}
}

type worker struct {
	ctx       context.Context
	converter converter.Converter
	mocker    Mocker
}

func (w *worker) GetType() common.WorkerKind {
	return common.WorkerKindMock
}

func (w *worker) Configure(ctx context.Context, arguments common.Arguments) error {
	templateContainers, err := w.converter.Convert(ctx, arguments)
	if err != nil {
		return err
	}

	return w.mocker.Configure(ctx, templateContainers)
}

func (w *worker) Run() error {
	return w.mocker.Run()
}
