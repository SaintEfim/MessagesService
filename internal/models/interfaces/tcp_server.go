package interfaces

import "context"

type TCPServer interface {
	AcceptConnection(ctx context.Context)
	RefuseConnection(ctx context.Context) error
}
