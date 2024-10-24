package server

import (
	"context"
	"net"
	"sync"

	"MessagesService/config"
	"MessagesService/internal/models/interfaces"

	"go.uber.org/zap"
)

type TCPServer struct {
	listener net.Listener
	handler  interfaces.MessageHandler
	logger   *zap.Logger
	cfg      *config.Config
	wg       *sync.WaitGroup
}

func NewWaitGroup() *sync.WaitGroup {
	return &sync.WaitGroup{}
}

func NewTCPListener(ctx context.Context, logger *zap.Logger, cfg *config.Config) (net.Listener, error) {
	listener, err := net.Listen(cfg.Server.Type, cfg.Server.Port)
	if err != nil {
		logger.Error("Error starting TCP server port" + cfg.Server.Port + err.Error())
		return nil, err
	}

	logger.Info("TCP server started port" + cfg.Server.Port)
	return listener, nil
}

func NewTCPServer(listener net.Listener, handler interfaces.MessageHandler, logger *zap.Logger, cfg *config.Config, wg *sync.WaitGroup) interfaces.TCPServer {
	return &TCPServer{
		listener: listener,
		handler:  handler,
		logger:   logger,
		cfg:      cfg,
		wg:       wg,
	}
}

func (s *TCPServer) AcceptLoop(ctx context.Context) error {
	for {
		errCh := make(chan error, 1)

		s.wg.Add(1)

		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting:" + err.Error())
			s.wg.Done()

			continue
		}

		go func(conn net.Conn) {
			defer s.wg.Done()
			defer conn.Close()

			s.logger.Info("Accepting connection from " + conn.RemoteAddr().String())

			if err := s.handler.MessageHandleRequest(ctx, conn); err != nil {
				s.logger.Error("Error handling request: " + err.Error())
				errCh <- err
			}
		}(conn)

		select {
		case <-ctx.Done():
			s.logger.Info("Closing TCP server listener")
			s.listener.Close()
			s.wg.Done()

			return ctx.Err()
		case handleErr := <-errCh:
			if handleErr != nil {
				s.logger.Error("Error handling request:" + handleErr.Error())
				s.listener.Close()
				s.wg.Wait()

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
