package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"MessagesService/config"
	"MessagesService/internal/models/entity"
	"MessagesService/internal/models/interfaces"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const keyForId = "nameid"

type MessageController struct {
	logger *zap.Logger
	cfg    *config.Config
	repo   interfaces.RedisRepository
}

func NewMessageController(logger *zap.Logger, cfg *config.Config, repo interfaces.RedisRepository) interfaces.MessageController {
	return &MessageController{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

func (c *MessageController) MessageProcessRequest(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(conn)

	msg, err := c.readTCPRequest(ctx, scanner, conn)
	if err != nil {
		c.logger.Error("Error processing TCP request:", zap.Error(err))
		return err
	}

	if len(msg.Message) == 0 {
		_, err := conn.Write([]byte("User init\n"))
		if err != nil {
			c.logger.Error("Error writing response to connection:", zap.Error(err))
			return err
		}
		c.logger.Info("Received empty message, responded with 'User init'.")
		return nil
	}

	if err := c.writeTCPRequest(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (c *MessageController) readTCPRequest(ctx context.Context, scanner *bufio.Scanner, conn net.Conn) (entity.TCPRequest, error) {
	var (
		msg entity.TCPRequest
	)

	for scanner.Scan() {
		clientMessage := scanner.Text()

		c.logger.Info("Received raw message", zap.String("clientMessage", clientMessage))

		if err := json.Unmarshal([]byte(clientMessage), &msg); err != nil {
			c.logger.Error("Error unmarshalling JSON", zap.Error(err))
			if _, err = conn.Write([]byte("Invalid JSON format.\n")); err != nil {
				c.logger.Error("Error sending response to client:", zap.Error(err))
			}
			return msg, err
		}

		userId, err := c.parseJWTToken(ctx, msg.UserCredential)
		if err != nil {
			c.logger.Error("Error parsing user ID", zap.Error(err))
			return msg, err
		}

		if err := c.repo.Set(ctx, userId.String(), conn.RemoteAddr().String()); err != nil {
			c.logger.Error("Error sending response to client:", zap.Error(err))
			return msg, err
		}

		c.logger.Info("Connection saved", zap.String("user_id", userId.String()))

		break
	}

	if err := scanner.Err(); err != nil {
		c.logger.Error("Error reading from connection", zap.Error(err))
		return msg, err
	}

	return msg, nil
}

func (c *MessageController) writeTCPRequest(ctx context.Context, message entity.TCPRequest) error {
	colleagueAddr, err := c.repo.Get(ctx, message.UserCredential.ColleagueId.String())
	if err != nil {
		c.logger.Error("Error getting colleague connection address:", zap.Error(err))
		return err
	}

	colleagueConn, err := net.Dial(c.cfg.Server.Type, colleagueAddr)
	if err != nil {
		c.logger.Error("Error connecting to colleague:", zap.Error(err))
		return err
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("Error marshalling message to JSON:", zap.Error(err))
		return err
	}

	_, err = colleagueConn.Write(messageJSON)
	if err != nil {
		c.logger.Error("Error sending message to colleague:", zap.Error(err))
		return err
	}

	c.logger.Info("Message sent to colleague", zap.String("colleague_addr", colleagueAddr))
	return nil
}

func (c *MessageController) parseJWTToken(ctx context.Context, user entity.UserCredential) (uuid.UUID, error) {
	var (
		claims = jwt.MapClaims{}
	)

	_, err := jwt.ParseWithClaims(user.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.cfg.AuthenticationConfiguration.AccessSecretKey), nil
	})
	if err != nil {
		c.logger.Error("Error parsing JWT token", zap.Error(err))
		return uuid.Nil, err
	}

	userIdStr, ok := claims[keyForId].(string)
	if !ok {
		c.logger.Error("user_id not found in token claims or invalid type")
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.logger.Error("Error parsing user ID", zap.Error(err))
		return uuid.Nil, err
	}

	return userId, nil
}
