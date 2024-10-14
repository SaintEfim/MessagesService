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
	defer func() {
		if err := conn.Close(); err != nil {
			c.logger.Error("Error closing connection:", zap.Error(err))
		}
	}()

	scanner := bufio.NewScanner(conn)

	msg, err := c.readTCPRequest(ctx, scanner, conn)
	if err != nil {
		c.logger.Error("Error processing TCP request:", zap.Error(err))
		return err
	}

	c.logger.Info("Successfully received message",
		zap.String("message", msg.Message),
		zap.String("colleague_id", msg.UserCredential.ColleagueId.String()),
	)

	return nil
}

func (c *MessageController) readTCPRequest(ctx context.Context, scanner *bufio.Scanner, conn net.Conn) (entity.TCPRequest, error) {
	var msg entity.TCPRequest

	for scanner.Scan() {
		clientMessage := scanner.Text()

		c.logger.Info("Received raw message", zap.String("clientMessage", clientMessage))

		err := json.Unmarshal([]byte(clientMessage), &msg)
		if err != nil {
			c.logger.Error("Error unmarshalling JSON", zap.Error(err))
			if _, err := conn.Write([]byte("Invalid JSON format.\n")); err != nil {
				c.logger.Error("Error sending response to client:", zap.Error(err))
				return msg, err
			}

			continue
		}

		if _, err := conn.Write([]byte("Message received.\n")); err != nil {
			c.logger.Error("Error sending response to client:", zap.Error(err))
			return msg, err
		}

		break
	}

	if err := scanner.Err(); err != nil {
		c.logger.Error("Error reading from connection", zap.Error(err))
		return msg, err
	}

	return msg, nil
}
