package interfaces

import (
	"context"
	"net"
)

type Controller interface {
	MessageHandleRequest(ctx context.Context, conn net.Conn) error
}
