package interfaces

import (
	"context"

	"MessagesService/proto/chat"
)

type ChatGrpcClient interface {
	Initialize(ctx context.Context) error
	Close(ctx context.Context) error
	chat.GreeterChatsClient
}
