package controller

import (
	"context"
	"log"
	"sync"
	"time"

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

func (c *Controller) SendMessage(ctx context.Context, req *dto.SendMessageRequest) (*dto.CreateActionResponse, error) {
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

func (c *Controller) handleSendMessage(ctx context.Context, req *dto.SendMessageRequest) (*dto.CreateActionResponse, error) {
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

	createdAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", resChatClient.CreateAt)
	if err != nil {
		log.Printf("Failed to parse CreateAt: %v", err)
		return nil, err
	}

	c.mu.Lock()
	receiver, exists := c.clients[req.ReceiverId]
	c.mu.Unlock()

	if exists {
		if err := receiver.TransferData(&dto.WsMessages{
			SenderId:   req.SenderId,
			ReceiverId: req.ReceiverId,
			Text:       req.Text,
			CreateAt:   createdAt,
		}); err != nil {
			c.mu.Lock()
			delete(c.clients, req.ReceiverId)
			c.mu.Unlock()
			return nil, err
		}
	}

	return &dto.CreateActionResponse{
		Id:        id,
		CreatedAt: createdAt,
	}, nil
}

func (c *Controller) Connect(ctx context.Context, client *dto.ConnectClientRequest, conn interfaces.Transfer) error {
	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		return err
	}

	c.mu.Lock()
	c.clients[client.Id] = conn
	c.mu.Unlock()

	return nil
}
