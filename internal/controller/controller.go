package controller

import (
	"context"
	"sync"
	"time"

	"MessagesService/config"
	"MessagesService/internal/models/dto"
	"MessagesService/internal/models/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Controller struct {
	logger  *zap.Logger
	cfg     *config.Config
	clients map[uuid.UUID]interfaces.Transfer
	mu      sync.Mutex
}

func NewController(logger *zap.Logger, cfg *config.Config) interfaces.Controller {
	return &Controller{
		logger:  logger,
		cfg:     cfg,
		clients: make(map[uuid.UUID]interfaces.Transfer),
		mu:      sync.Mutex{},
	}
}

func (c *Controller) SendMessage(ctx context.Context, req *dto.SendMessage, conn interfaces.Transfer) error {
	validate := validator.New()
	conSender := c.clients[req.SenderId]

	if err := validate.Struct(req); err != nil {
		if err := conSender.TransferData(&dto.ResponseMessage{Error: "Failed validate: " + err.Error()}); err != nil {
			return err
		}
	}

	return c.handleSendMessage(ctx, req)
}

func (c *Controller) handleSendMessage(ctx context.Context, req *dto.SendMessage) error {
	msg := &dto.ResponseMessage{
		Text:      req.Text,
		CreatedAt: time.Now(),
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	receiver, exists := c.clients[req.ReceiverId]
	if !exists {
		c.logger.Warn("Receiver not found",
			zap.String("receiver_id", req.ReceiverId.String()))
		// TODO: Сохранить сообщение в БД через gRPC
		return nil
	}

	if err := receiver.TransferData(msg); err != nil {
		c.logger.Error("Failed to send message",
			zap.String("receiver_id", req.ReceiverId.String()),
			zap.Error(err))

		delete(c.clients, req.ReceiverId)
		return err
	}

	c.logger.Info("Successfully send message: " + req.Text)

	// TODO: Сохранить сообщение в БД после успешной отправки
	return nil
}

func (c *Controller) Connect(ctx context.Context, client *dto.ConnectClient, conn interfaces.Transfer) error {
	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		_ = conn.TransferData(&dto.ResponseMessage{Error: "Failed validate: " + err.Error()})
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.clients[client.Id] = conn

	return nil
}
