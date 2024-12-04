package main

import (
	"context"

	"MessagesService/config"
	"MessagesService/internal/controller"
	"MessagesService/internal/handler"
	"MessagesService/internal/models/interfaces"
	"MessagesService/internal/repository/redis"
	"MessagesService/internal/server"
	"MessagesService/pkg/logger"

	redisClient "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func registerRedis(lc fx.Lifecycle, redisClient *redisClient.Client) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := redisClient.Ping(ctx).Err(); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := redisClient.Close(); err != nil {
				return err
			}
			return nil
		},
	})
}

func registerServer(lifecycle fx.Lifecycle,
	mainCtx context.Context,
	srv interfaces.TCPServer,
	logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.AcceptConnection(mainCtx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server...")

			if err := srv.RefuseConnection(mainCtx); err != nil {
				logger.Error("Failed to stop server" + err.Error())
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
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			logger.NewLogger,
			redis.NewRedisClient,
			redis.NewRedisRepository,
			controller.NewController,
			handler.NewHandler,
			server.NewTCPServer,
		),
		fx.Invoke(registerServer),
		fx.Invoke(registerRedis),
	).Run()
}
