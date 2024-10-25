package interfaces

import "context"

type TCPServer interface {
	AcceptConnection(ctx context.Context) error
	RefuseConnection(ctx context.Context) error
}
