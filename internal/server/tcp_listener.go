package server

import (
	"context"
	"net"

	"MessagesService/config"

	"go.uber.org/zap"
)

func NewTCPListener(ctx context.Context, cfg *config.Config, logger *zap.Logger) (net.Listener, error) {
	listener, err := net.Listen(cfg.Server.Type, cfg.Server.Port)
	if err != nil {
		logger.Error("Error starting TCP server", zap.String("port", cfg.Server.Port), zap.Error(err))
		return nil, err
	}

	logger.Info("TCP server started", zap.String("port", cfg.Server.Port))
	return listener, nil
}
