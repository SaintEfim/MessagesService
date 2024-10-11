package handler

import (
	"context"
	"go.uber.org/zap"
	"net"

	"MessageService/internal/models/interfaces"
)

type Handler struct {
	controller interfaces.Controller
	logger     *zap.Logger
}

func NewHandler(controller interfaces.Controller, logger *zap.Logger) interfaces.Handler {
	return &Handler{
		controller: controller,
		logger:     logger,
	}
}

func (h *Handler) HandleConnection(ctx context.Context, conn net.Conn) error {
	var err error
	err = h.controller.Connection(ctx, conn)

	if err != nil {
		h.logger.Error("Error connection:", zap.Error(err))
		return err
	}

	return nil
}
