package interfaces

import (
	"context"

	"MessagesService/internal/models/dto"
)

type Controller interface {
	SendMessage(ctx context.Context, req *dto.SendMessage) (*dto.CreateAction, error)
	Connect(ctx context.Context, client *dto.ConnectClient, conn Transfer) error
}
