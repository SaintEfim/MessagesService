package server

import (
	"MessagesService/config"
	"MessagesService/internal/models/interfaces"
	"context"
	"net"

	"go.uber.org/zap"
)

type TCPServer struct {
	listener net.Listener
	handler  interfaces.MessageHandler
	logger   *zap.Logger
	errCh    chan error
	cfg      *config.Config
}

func NewTCPListener(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config) (net.Listener, error) {
	listener, err := net.Listen(cfg.Server.Type, cfg.Server.Port)
	if err != nil {
		logger.Error("Error starting TCP server port" + cfg.Server.Port + err.Error())
		return nil, err
	}

	logger.Info("TCP server started port" + cfg.Server.Port)
	return listener, nil
}

func NewTCPServer(
	listener net.Listener,
	handler interfaces.MessageHandler,
	logger *zap.Logger,
	errCh chan error,
	cfg *config.Config) interfaces.TCPServer {
	return &TCPServer{
		listener: listener,
		handler:  handler,
		logger:   logger,
		errCh:    errCh,
		cfg:      cfg,
	}
}

func (s *TCPServer) AcceptConnection(ctx context.Context) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic("Error accepting:" + err.Error())
		}

		go func(conn net.Conn, ctx context.Context) {
			defer conn.Close()

			s.logger.Info("Accepting connection from " + conn.RemoteAddr().String())

			if err := s.handler.MessageHandleRequest(ctx, conn); err != nil {
				s.logger.Error("Error handling request: " + err.Error())
				s.errCh <- err
			}

			select {
			case err := <-s.errCh:
				s.logger.Error("не знаю что делать......." + err.Error())
			}
		}(conn, ctx)
	}
}

func (s *TCPServer) RefuseConnection(ctx context.Context) error {
	defer close(s.errCh)

	if err := s.listener.Close(); err != nil {
		s.logger.Error("Error closing:" + err.Error())
		return err
	}
	return nil
}
