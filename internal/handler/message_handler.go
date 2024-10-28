package handler

import (
	"context"
	"net"

	"MessagesService/internal/models/interfaces"

	"go.uber.org/zap"
)

type MessageHandler struct {
	controller interfaces.MessageController
	logger     *zap.Logger
}

func NewMessageHandler(controller interfaces.MessageController, logger *zap.Logger) interfaces.MessageHandler {
	return &MessageHandler{
		controller: controller,
		logger:     logger,
	}
}

func (h *MessageHandler) MessageHandleRequest(ctx context.Context, conn net.Conn) error {
	var err = h.controller.MessageHandleRequest(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}
