package controller

import (
	"context"

	"MessagesService/config"
	"MessagesService/internal/models/dto"
	"MessagesService/internal/models/interfaces"

	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
	cfg    *config.Config
}

func NewController(logger *zap.Logger, cfg *config.Config) interfaces.Controller {
	return &Controller{
		logger: logger,
		cfg:    cfg,
	}
}

func (c Controller) SendMessage(ctx context.Context, req *dto.SendMessage, conn interfaces.Transfer) error {
	//TODO implement me
	panic("implement me")
}
