package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"MessagesService/internal/models/entity"
	"MessagesService/internal/models/interfaces"

	"go.uber.org/zap"
)

type MessageController struct {
	logger *zap.Logger
}

func NewMessageController(logger *zap.Logger) interfaces.MessageController {
	return &MessageController{
		logger: logger,
	}
}

func (c *MessageController) MessageProcessRequest(ctx context.Context, conn net.Conn) error {
	var (
		err error
		msg entity.TCPRequest
	)
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

		err := json.Unmarshal([]byte(clientMessage), &msg)
		if err != nil {
			c.logger.Error("Error parsing JSON", zap.Error(err))
			conn.Write([]byte("Invalid JSON format.\n"))
			if err != nil {
				c.logger.Error("Error sending response:", zap.Error(err))
				return err
			}
			continue
		}

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
