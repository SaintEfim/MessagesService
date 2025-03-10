package controller

import (
	"context"
	"sync"

	"MessagesService/config"
	"MessagesService/internal/models/dto"
	"MessagesService/internal/models/interfaces"
	"MessagesService/proto/chat"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Controller struct {
	logger     *zap.Logger
	cfg        *config.Config
	clients    map[uuid.UUID]interfaces.Transfer
	chatClient interfaces.ChatGrpcClient
	mu         sync.Mutex
}

func NewController(logger *zap.Logger, cfg *config.Config, chatClient interfaces.ChatGrpcClient) interfaces.Controller {
	return &Controller{
		logger:     logger,
		cfg:        cfg,
		clients:    make(map[uuid.UUID]interfaces.Transfer),
		chatClient: chatClient,
		mu:         sync.Mutex{},
	}
}

func (c *Controller) SendMessage(ctx context.Context, req *dto.SendMessage) (*dto.CreateAction, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	res, err := c.handleSendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) handleSendMessage(ctx context.Context, req *dto.SendMessage) (*dto.CreateAction, error) {
	resChatClient, err := c.chatClient.CreateMessage(ctx, &chat.MessageCreateRequest{
		ChatId:     req.ChatId.String(),
		SenderId:   req.SenderId.String(),
		ReceiverId: req.ReceiverId.String(),
		Text:       req.Text,
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(resChatClient.Id)
	if err != nil {
		return nil, err
	}

	res := &dto.CreateAction{
		Id: id,
	}

	reqReceiver := &dto.Messages{
		Text: req.Text,
	}

	c.mu.Lock()
	receiver, exists := c.clients[req.ReceiverId]
	c.mu.Unlock()

	if exists {
		if err := receiver.TransferData(reqReceiver); err != nil {
			c.mu.Lock()
			delete(c.clients, req.ReceiverId)
			c.mu.Unlock()
			return nil, err
		}
	}

	return res, nil
}

func (c *Controller) Connect(ctx context.Context, client *dto.ConnectClient, conn interfaces.Transfer) error {
	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.clients[client.Id] = conn

	return nil
}
