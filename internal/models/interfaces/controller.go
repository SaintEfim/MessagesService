package interfaces

import (
	"context"

	"MessagesService/internal/models/dto"
)

type Controller interface {
	SendMessage(ctx context.Context, req *dto.SendMessageRequest) (*dto.CreateActionResponse, error)
	Connect(ctx context.Context, client *dto.ConnectClientRequest, conn Transfer) error
}
