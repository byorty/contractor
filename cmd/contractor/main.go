package main

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"github.com/byorty/contractor/logger"
	"github.com/byorty/contractor/mocker"
	"github.com/byorty/contractor/tester"
)

func main() {
	app := common.NewApplication(
		common.Constructors,
		converter.Constructors,
		logger.Constructors,
		tester.Constructors,
		mocker.Constructors,
	)
	app.Run(func(
		ctx context.Context,
		args common.Arguments,
		container common.WorkerContainer,
	) error {
		worker, err := container.Get(args.Mode)
		if err != nil {
			return err
		}

		err = worker.Configure(ctx, args)
		if err != nil {
			return err
		}

		return worker.Run()
	})
}
