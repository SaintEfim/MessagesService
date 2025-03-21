package main

import (
	"context"

	"MessagesService/config"
	"MessagesService/internal/controller"
	grpcSrv "MessagesService/internal/delivery/grpc"
	websocketHandler "MessagesService/internal/handler/websocket"
	"MessagesService/internal/models/interfaces"
	websocketSrv "MessagesService/internal/server/websocket"
	"MessagesService/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func registerServer(lc fx.Lifecycle, srv interfaces.Server) {
	lc.Append(fx.Hook{
		OnStart: srv.Run,
		OnStop:  srv.Stop,
	})
}

func registerGRPCClient(lc fx.Lifecycle, client interfaces.ChatGrpcClient) {
	lc.Append(fx.Hook{
		OnStart: client.Initialize,
		OnStop:  client.Close,
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
			controller.NewController,
			grpcSrv.NewChatGrpcClient,
			websocketHandler.NewHandler,
			websocketSrv.NewServer,
			websocketSrv.NewUpgrader,
		),
		fx.Invoke(registerServer),
		fx.Invoke(registerGRPCClient),
	).Run()
}
