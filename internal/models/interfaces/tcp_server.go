package interfaces

import "context"

type TCPServer interface {
	AcceptLoop(ctx context.Context) error
	RefuseLoop(ctx context.Context) error
}
