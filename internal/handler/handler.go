package handler

import (
	"context"
	"net"

	"MessagesService/internal/models/interfaces"

	"go.uber.org/zap"
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

func (h *Handler) MessageHandleRequest(ctx context.Context, conn net.Conn) error {
	var err = h.controller.MessageHandleRequest(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}
