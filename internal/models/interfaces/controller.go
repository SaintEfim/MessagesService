package interfaces

import (
	"context"

	"MessagesService/internal/models/dto"
)

type Controller interface {
	SendMessage(ctx context.Context, req *dto.SendMessage, conn Transfer) error
}
