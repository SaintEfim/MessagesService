package server

import (
	"context"
	"net"

	"MessagesService/config"
	"MessagesService/internal/models/interfaces"

	"go.uber.org/zap"
)

type TCPServer struct {
	listener net.Listener
	cfg      *config.Config
	handler  interfaces.MessageHandler
	logger   *zap.Logger
}

func NewTCPListener(ctx context.Context, cfg *config.Config, logger *zap.Logger) (net.Listener, error) {
	listener, err := net.Listen(cfg.Server.Type, cfg.Server.Port)
	if err != nil {
		logger.Error("Error starting TCP server", zap.String("port", cfg.Server.Port), zap.Error(err))
		return nil, err
	}

	logger.Info("TCP server started", zap.String("port", cfg.Server.Port))
	return listener, nil
}

func NewTCPServer(listener net.Listener, cfg *config.Config, handler interfaces.MessageHandler, logger *zap.Logger) interfaces.TCPServer {
	return &TCPServer{
		listener: listener,
		cfg:      cfg,
		handler:  handler,
		logger:   logger,
	}
}

func (s *TCPServer) AcceptLoop(ctx context.Context) error {
	for {
		var conn, err = s.listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting:", zap.Error(err))
			return err
		}

		s.logger.Info("Connected with", zap.String("address", conn.RemoteAddr().String()))

		var handleErr error

		go func() {
			if err := s.handler.MessageHandleRequest(ctx, conn); err != nil {
				s.logger.Error("Error handling request:", zap.Error(err))
				handleErr = err
			}
		}()

		if handleErr != nil {
			return handleErr
		}
	}
}

func (s *TCPServer) RefuseLoop(ctx context.Context) error {
	if err := s.listener.Close(); err != nil {
		s.logger.Error("Error closing:", zap.Error(err))
	}
	return nil
}
