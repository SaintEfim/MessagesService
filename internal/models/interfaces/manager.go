package interfaces

import "context"

type Manager interface {
	Start(ctx context.Context) error
}
