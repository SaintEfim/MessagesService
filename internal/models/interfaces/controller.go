package interfaces

import (
	"context"
	"net"
)

type Controller interface {
	Connection(ctx context.Context, conn net.Conn) error
}
