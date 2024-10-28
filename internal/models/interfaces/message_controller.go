package interfaces

import (
	"context"
	"net"
)

type MessageController interface {
	MessageHandleRequest(ctx context.Context, conn net.Conn) error
}
