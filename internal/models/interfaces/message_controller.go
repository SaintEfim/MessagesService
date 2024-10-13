package interfaces

import (
	"context"
	"net"
)

type MessageController interface {
	MessageProcessRequest(ctx context.Context, conn net.Conn) error
}
