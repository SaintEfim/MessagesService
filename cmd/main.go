package main

import (
	"context"

	"MessageService/config"
	"MessageService/internal/controller"
	"MessageService/internal/handler"
	"MessageService/internal/models/interfaces"
	"MessageService/internal/server"
	"MessageService/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func registerServer(lifecycle fx.Lifecycle, srv interfaces.Server, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server...")

			var err error
			go func() {
				err = srv.Run(ctx)
				if err != nil {
					logger.Error("Server failed to start", zap.Error(err))
				} else {
					logger.Info("Server started successfully")
				}
			}()

			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server...")

			if err := srv.Stop(ctx); err != nil {
				logger.Error("Failed to stop server", zap.Error(err))
				return err
			}

			logger.Info("Server stopped successfully")
			return nil
		},
	})
}

func main() {
	fx.New(
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		fx.Provide(func() (*config.Config, error) {
			return config.ReadConfig("config", "yaml", "./config")
		}),
		fx.Provide(
			logger.NewLogger,
			server.NewTCPServer,
			server.NewServer,
			handler.NewHandler,
			controller.NewController),
		fx.Invoke(registerServer),
	).Run()
}
