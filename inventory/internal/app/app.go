package app

import (
	"context"
	"fmt"
	"net"

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

func (a *App) Run(ctx context.Context) error {
	cfg, err := a.serviceProvider.Config()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(cfg.GRPC().Host(), fmt.Sprintf("%d", cfg.GRPC().Port()))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	server := a.serviceProvider.GRPCServer()
	cl := a.serviceProvider.Closer()

	cl.AddNamed("gRPC server", func(_ context.Context) error {
		server.GracefulStop()
		return nil
	})
	cl.AddNamed("listener", func(_ context.Context) error {
		return lis.Close()
	})

	errCh := make(chan error, 1)
	go func() {
		logger.Info(ctx, fmt.Sprintf("🚀 Inventory gRPC server listening on %s", addr))

		if err := server.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "🛑 Получен сигнал завершения, останавливаем Inventory...")

		return cl.CloseAll(context.WithoutCancel(ctx))
	case err := <-errCh:
		return fmt.Errorf("failed to serve Inventory gRPC server: %w", err)
	}
}
