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
	handler  interfaces.Handler
	logger   *zap.Logger
	cfg      *config.Config
}

func NewTCPServer(
	ctx context.Context,
	handler interfaces.Handler,
	logger *zap.Logger,
	cfg *config.Config) (interfaces.TCPServer, error) {
	listenerTCP, err := net.Listen(cfg.Server.Type, cfg.Server.Port)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		listener: listenerTCP,
		handler:  handler,
		logger:   logger,
		cfg:      cfg,
	}, nil
}

func (s *TCPServer) AcceptConnection(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting connection" + err.Error())
			continue
		}

		go func(conn net.Conn, ctx context.Context) error {
			defer conn.Close()

			s.logger.Info("Accepting connection from " + conn.RemoteAddr().String())

			if err := s.handler.MessageHandleRequest(ctx, conn); err != nil {
				s.logger.Error("Error handling request: " + err.Error())
				return err
			}

			return nil
		}(conn, ctx)
	}
}

func (s *TCPServer) RefuseConnection(ctx context.Context) error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}
