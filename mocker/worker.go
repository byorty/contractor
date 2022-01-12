package mocker

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"go.uber.org/fx"
)

type WorkerIn struct {
	fx.In
	Ctx                context.Context
	ConverterContainer converter.Container
	Mocker             Mocker
}

func NewFxWorker(in WorkerIn) common.Worker {
	return &worker{
		ctx:                in.Ctx,
		converterContainer: in.ConverterContainer,
		mocker:             in.Mocker,
	}
}

type worker struct {
	ctx                context.Context
	converterContainer converter.Container
	mocker             Mocker
}

func (w *worker) GetType() common.WorkerType {
	return common.WorkerTypeMock
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

	return w.mocker.Configure(ctx, templateContainers)
}

func (w *worker) Run() error {
	return w.mocker.Run()
}
