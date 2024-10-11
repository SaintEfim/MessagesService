package interfaces

import (
	"context"
	"net"
)

type Handler interface {
	HandleConnection(ctx context.Context, conn net.Conn) error
}
