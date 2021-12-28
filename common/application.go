package common

import (
	"context"
	"github.com/jessevdk/go-flags"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	ctx     context.Context
	cancel  context.CancelFunc
	fxApp   *fx.App
	options []fx.Option
}

func NewApplication(providers ...interface{}) *Application {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	app := &Application{
		options: []fx.Option{
			fx.NopLogger,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	for _, provider := range providers {
		switch p := provider.(type) {
		case fx.Option:
			app.options = append(app.options, p)
		default:
			app.options = append(app.options, fx.Provide(p))
		}
	}

	return app
}

func (a *Application) Run(invoker interface{}) {
	var args Arguments
	_, err := flags.Parse(&args)
	if err != nil {
		panic(err)
	}

	a.options = append(
		a.options,
		fx.Provide(func() context.Context {
			return a.ctx
		}),
		fx.Provide(func() Arguments {
			return args
		}),
		fx.Invoke(invoker),
	)
	a.fxApp = fx.New(a.options...)

	go a.listenSignals()

	startCtx, cancel := context.WithTimeout(a.ctx, fx.DefaultTimeout)
	defer cancel()

	if err = a.fxApp.Start(startCtx); err != nil {
		panic(err)
	}
}

func (a *Application) listenSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//for sig := range signals {
	//	a.logger.Infof("income signal %s", sig)
	//	a.Stop()
	//	a.cancel()
	//	return
	//}
}
