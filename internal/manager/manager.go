package manager

import (
	"MessagesService/internal/models/interfaces"
	"context"
)

type Manager struct {
	srv   interfaces.TCPServer
	errCh chan error
}

func NewManager(srv interfaces.TCPServer, errCh chan error) interfaces.Manager {
	return &Manager{
		srv:   srv,
		errCh: errCh,
	}
}

func (m *Manager) Start(ctx context.Context) error {
	go func() {
		m.srv.AcceptConnection(ctx, m.errCh)
	}()

	select {
	case err := <-m.errCh:
		return err
	}

	return nil
}
