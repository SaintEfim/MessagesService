package interfaces

import (
	"context"
	"net"
)

type MessageHandler interface {
	MessageHandleRequest(ctx context.Context, conn net.Conn) error
}
