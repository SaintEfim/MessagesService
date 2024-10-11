package controller

import (
	"MessageService/internal/models/interfaces"
	"bufio"
	"context"
	"go.uber.org/zap"
	"net"
)

type Controller struct {
	logger *zap.Logger
}

func NewController(logger *zap.Logger) interfaces.Controller {
	return &Controller{
		logger: logger,
	}
}

func (c *Controller) Connection(ctx context.Context, conn net.Conn) error {
	var err error
	defer func(conn net.Conn) {
		err = conn.Close()
	}(conn)

	if err != nil {
		c.logger.Error("Error closing:", zap.Error(err))
		return err
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		clientMessage := scanner.Text()

		c.logger.Info("Received from client", zap.String("clientMessage", clientMessage))
		_, err = conn.Write([]byte("Message received.\n"))
		if err != nil {
			c.logger.Error("Error sending response:", zap.Error(err))
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		c.logger.Error("Error reading:", zap.Error(err))
		return err
	}

	return nil
}
