package common

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var (
	ErrUnsupportedWorkerType = errors.New("unsupported worker type")
)

type WorkerType string

const (
	WorkerTypeMock WorkerType = "mock"
	WorkerTypeTest WorkerType = "test"
	WorkerTypeE2E  WorkerType = "e2e"
)

type Worker interface {
	GetType() WorkerType
	Configure(ctx context.Context, arguments Arguments) error
	Run() error
}

type WorkerContainerFxIn struct {
	fx.In
	Workers []Worker `group:"worker"`
}

func NewFxWorkerContainer(in WorkerContainerFxIn) WorkerContainer {
	container := make(WorkerContainer)

	for _, worker := range in.Workers {
		container[worker.GetType()] = worker
	}

	return container
}

type WorkerContainer map[WorkerType]Worker

func (c WorkerContainer) Get(workerType WorkerType) (Worker, error) {
	worker, ok := c[workerType]
	if !ok {
		return nil, errors.Wrap(ErrUnsupportedWorkerType, string(workerType))
	}

	return worker, nil
}
