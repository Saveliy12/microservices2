package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{
		serviceProvider: newServiceProvider(),
	}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	_, err := a.serviceProvider.Config()

	return err
}

func (a *App) initLogger(_ context.Context) error {
	cfg, err := a.serviceProvider.Config()
	if err != nil {
		return err
	}

	loggerCfg := cfg.Logger()

	if err := logger.Init(loggerCfg.Level(), loggerCfg.AsJSON()); err != nil {
		return err
	}

	a.serviceProvider.Closer().AddNamed("logger", func(_ context.Context) error {
		return logger.Sync()
	})

	return nil
}

// Run запускает HTTP сервер и блокируется до ctx.Done, затем graceful shutdown через closer
func (a *App) Run(ctx context.Context) error {
	server, err := a.serviceProvider.HTTPServer(ctx)
	if err != nil {
		return err
	}

	cl := a.serviceProvider.Closer()

	cl.AddNamed("http server", func(c context.Context) error {
		return server.Shutdown(c)
	})

	errCh := make(chan error, 1)
	go func() {
		logger.Info(ctx, fmt.Sprintf("🚀 Order HTTP server listening on %s", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "🛑 Получен сигнал завершения, останавливаем Order...")

		return cl.CloseAll(context.WithoutCancel(ctx))
	case err := <-errCh:
		return fmt.Errorf("failed to serve Order HTTP server: %w", err)
	}
}
