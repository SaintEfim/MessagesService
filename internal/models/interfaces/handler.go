package interfaces

import (
	"context"
	"net"
)

type Handler interface {
	MessageHandleRequest(ctx context.Context, conn net.Conn) error
}
