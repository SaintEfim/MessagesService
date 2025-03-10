package grpc

import (
	"context"

	"MessagesService/config"
	"MessagesService/internal/models/interfaces"
	"MessagesService/proto/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatGrpcClient struct {
	client chat.GreeterChatsClient
	cfg    *config.Config
	conn   *grpc.ClientConn
}

func NewChatGrpcClient(ctx context.Context, cfg *config.Config) interfaces.ChatGrpcClient {
	return &ChatGrpcClient{
		cfg: cfg,
	}
}

func (c *ChatGrpcClient) Initialize(ctx context.Context) error {
	conn, err := grpc.NewClient(c.cfg.GRPCClient.Services["chats"], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = chat.NewGreeterChatsClient(c.conn)

	return nil
}

func (c *ChatGrpcClient) Close(ctx context.Context) error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *ChatGrpcClient) CreateMessage(ctx context.Context, in *chat.MessageCreateRequest, opts ...grpc.CallOption) (*chat.MessageCreateResponse, error) {
	res, err := c.client.CreateMessage(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}
