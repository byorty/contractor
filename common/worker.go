package common

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
)

type WorkerKind string

const (
	WorkerKindMock WorkerKind = "mock"
	WorkerKindTest WorkerKind = "test"
)

type Worker interface {
	GetType() WorkerKind
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

type WorkerContainer map[WorkerKind]Worker

func (c WorkerContainer) Get(kind WorkerKind) (Worker, error) {
	worker, ok := c[kind]
	if !ok {
		return nil, errors.New(fmt.Sprintf("worket is not found by kind=%s", kind))
	}

	return worker, nil
}
