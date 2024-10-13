package handler

import (
	"context"
	"go.uber.org/zap"
	"net"

	"MessagesService/internal/models/interfaces"
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
	var err error
	err = h.controller.MessageProcessRequest(ctx, conn)

	if err != nil {
		h.logger.Error("Error connection:", zap.Error(err))
		return err
	}

	return nil
}
