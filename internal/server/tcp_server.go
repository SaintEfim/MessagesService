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
		logger.Error("Error starting TCP server port" + cfg.Server.Port + err.Error())
		return nil, err
	}

	logger.Info("TCP server started port" + cfg.Server.Port)
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
		errCh := make(chan error, 1)
		conn, err := s.listener.Accept()

		if err != nil {
			s.logger.Error("Error accepting:" + err.Error())
			return err
		}

		go func() {
			if err := s.handler.MessageHandleRequest(ctx, conn); err != nil {
				s.logger.Error("Error handling request:" + err.Error())
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			s.logger.Info("Closing TCP server listener")
			return ctx.Err()
		case handleErr := <-errCh:
			if handleErr != nil {
				s.logger.Error("Error handling request:" + handleErr.Error())
				return err
			}
		default:
			s.logger.Debug("Waiting for connection or error")
		}
	}
}

func (s *TCPServer) RefuseLoop(ctx context.Context) error {
	if err := s.listener.Close(); err != nil {
		s.logger.Error("Error closing:" + err.Error())
		return err
	}
	return nil
}
