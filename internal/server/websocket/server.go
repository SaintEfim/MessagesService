package websocket

import (
	"context"
	"errors"
	"net"
	"net/http"

	"MessagesService/config"
	"MessagesService/internal/middleware"
	"MessagesService/internal/models/interfaces"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	handler  interfaces.Handler
	logger   *zap.Logger
	srv      *http.Server
	cfg      *config.Config
	upgrader *websocket.Upgrader
}

func NewUpgrader(cfg *config.Config) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowedOrigins := map[string]bool{}
			for _, value := range cfg.Cors.AllowedOrigins {
				allowedOrigins[value] = true
			}
			return allowedOrigins[origin]
		},
	}
}

func NewServer(handler interfaces.Handler, logger *zap.Logger, cfg *config.Config, upgrader *websocket.Upgrader) interfaces.Server {
	httpServer := &http.Server{
		Addr: net.JoinHostPort(cfg.Server.Addr, cfg.Server.Port),
	}

	return &Server{
		handler:  handler,
		logger:   logger,
		srv:      httpServer,
		cfg:      cfg,
		upgrader: upgrader,
	}
}

func (s *Server) Run(ctx context.Context) error {
	var err error
	go func() {
		r := mux.NewRouter()
		r.Use(middleware.AuthMiddleware(s.logger, s.cfg.AuthenticationConfiguration.AccessSecretKey))

		handler := CorsSettings(s.cfg).Handler(r)
		s.srv.Handler = handler

		if err = s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("failed to start server", zap.Error(err))
		}
	}()

	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Server.Timeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
