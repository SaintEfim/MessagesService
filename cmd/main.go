package main

import (
	"context"

	"MessagesService/config"
	"MessagesService/internal/controller"
	"MessagesService/internal/handler"
	"MessagesService/internal/models/interfaces"
	"MessagesService/internal/server"
	"MessagesService/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func registerServer(lifecycle fx.Lifecycle, srv interfaces.TCPServer, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server...")

			var err error
			go func() {
				err = srv.AcceptLoop(ctx)
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

			if err := srv.RefuseLoop(ctx); err != nil {
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
			server.NewTCPListener,
			server.NewTCPServer,
			handler.NewMessageHandler,
			controller.NewMessageController),
		fx.Invoke(registerServer),
	).Run()
}
