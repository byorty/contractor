package main

import (
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"github.com/byorty/contractor/mocker"
	"github.com/byorty/contractor/tester"
	"github.com/byorty/contractor/tester/assertion"
	"github.com/byorty/contractor/tester/graylog"
	"github.com/byorty/contractor/tester/reporter"
)

func main() {
	app := common.NewApplication(
		common.Constructors,
		converter.Constructors,
		tester.Constructors,
		mocker.Constructors,
		assertion.Constructors,
		graylog.Constructors,
		reporter.Constructors,
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
			fmt.Println(err)
			return err
		}

		return worker.Run()
	})
}
