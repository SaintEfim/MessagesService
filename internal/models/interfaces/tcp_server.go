package interfaces

import "context"

type TCPServer interface {
	AcceptConnection(ctx context.Context, errCh chan error) error
	RefuseConnection(ctx context.Context) error
}
