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
	"go.uber.org/zap"
)

func registerRedis(lc fx.Lifecycle, mainCtx context.Context, redisClient *redisClient.Client, cfg *config.Config) {
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

func registerServer(lifecycle fx.Lifecycle, mainCtx context.Context, srv interfaces.TCPServer, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server...")

			errCh := make(chan error, 1)

			go func() {
				err := srv.AcceptLoop(mainCtx)
				if err != nil {
					logger.Error("Server failed to start: " + err.Error())
					errCh <- err
				} else {
					logger.Info("Server started successfully")
				}
			}()

			select {
			case <-ctx.Done():
				logger.Info("Context cancelled, stopping server")
				return ctx.Err()
			case err := <-errCh:
				if err != nil {
					return err
				}
			default:
				logger.Info("Server started successfully")
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server...")

			if err := srv.RefuseLoop(mainCtx); err != nil {
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
		fx.Provide(
			logger.NewLogger,
			redis.NewRedisClient,
			redis.NewRedisRepository,
			server.NewTCPListener,
			server.NewTCPServer,
			handler.NewMessageHandler,
			controller.NewMessageController),
		fx.Invoke(registerServer),
		fx.Invoke(registerRedis),
	).Run()
}
