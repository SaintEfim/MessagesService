package server

import (
	"context"
	"net"

	"MessageService/config"
	"MessageService/internal/models/interfaces"

	"go.uber.org/zap"
)

type Server struct {
	listener net.Listener
	cfg      *config.Config
	handler  interfaces.Handler
	logger   *zap.Logger
}

func NewServer(listener net.Listener, cfg *config.Config, handler interfaces.Handler, logger *zap.Logger) interfaces.Server {
	return &Server{
		listener: listener,
		cfg:      cfg,
		handler:  handler,
		logger:   logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	for {
		var err error
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting:", zap.Error(err))
		}

		s.logger.Info("Connected with", zap.String("address", conn.RemoteAddr().String()))

		go func() {
			err = s.handler.HandleConnection(ctx, conn)
		}()

		if err != nil {
			return err
		}
	}
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.listener.Close()
	if err != nil {
		s.logger.Error("Error closing:", zap.Error(err))
	}
	return nil
}
