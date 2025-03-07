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
	if err := validate.Struct(req); err != nil {
		if err := conn.TransferDataText("Failed validate: " + err.Error()); err != nil {
			return err
		}

		return err
	}

	return c.receiveMessage(ctx, req, conn)
}

func (c *Controller) Connect(ctx context.Context, client *dto.Connect, conn interfaces.Transfer) error {
	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		if err := conn.TransferDataText("Failed validate: " + err.Error()); err != nil {
			return err
		}

		return err
	}

	c.mu.Lock()
	if _, exists := c.clients[client.Id]; !exists {
		c.clients[client.Id] = conn
	}
	c.mu.Unlock()

	return nil
}

func (c *Controller) receiveMessage(ctx context.Context, req *dto.SendMessage, conn interfaces.Transfer) error {
	msg := &dto.ReceiveMessage{
		Content:   req.Message,
		CreatedAt: time.Now(),
	}

	c.mu.Lock()
	if _, exists := c.clients[req.SenderID]; !exists {
		c.clients[req.SenderID] = conn
	}
	c.mu.Unlock()

	c.mu.Lock()
	receiver, exists := c.clients[req.ReceiverID]
	c.mu.Unlock()

	if !exists {
		// TODO: Сохранить сообщение в БД через gRPC
	} else if err := receiver.TransferData(msg); err != nil {
		return err
	}

	// TODO: Сохранить сообщение в БД после успешной отправки
	return nil
}
