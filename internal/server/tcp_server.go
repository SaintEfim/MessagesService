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
		}
		s.logger.Info("Connected with", zap.String("address", conn.RemoteAddr().String()))

		go func() {
			if err := s.handler.MessageHandleRequest(ctx, conn); err != nil {
				s.logger.Error("Error handling request:", zap.Error(err))
			}
		}()

		if err != nil {
			return err
		}
	}
}

func (s *TCPServer) RefuseLoop(ctx context.Context) error {
	if err := s.listener.Close(); err != nil {
		s.logger.Error("Error closing:", zap.Error(err))
	}
	return nil
}
