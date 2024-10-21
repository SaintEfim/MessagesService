package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"net"

	"MessagesService/config"
	"MessagesService/internal/models/entity"
	"MessagesService/internal/models/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessageController struct {
	logger *zap.Logger
	cfg    *config.Config
	repo   interfaces.CacheRepository
}

func NewMessageController(logger *zap.Logger, cfg *config.Config, repo interfaces.CacheRepository) interfaces.MessageController {
	return &MessageController{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

func (c *MessageController) MessageProcessRequest(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(conn)

	for {
		msg, err := c.readTCPRequest(ctx, scanner, conn)
		if err != nil {
			_, err := conn.Write([]byte("Read message Error: " + err.Error() + "\n"))
			if err != nil {
				return err
			}

			continue
		}

		switch msg.Operation {
		case entity.OperationInit:
			continue
		case entity.OperationSendMessage:
			if err := c.writeTCPRequest(ctx, msg); err != nil {
				return err
			}
		default:
			_, err := conn.Write([]byte("Operation not found\n"))
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (c *MessageController) readTCPRequest(ctx context.Context, scanner *bufio.Scanner, conn net.Conn) (*entity.TCPRequest, error) {
	msg := entity.TCPRequest{}
	validate := validator.New()

	for scanner.Scan() {
		clientMessage := scanner.Text()

		if err := json.Unmarshal([]byte(clientMessage), &msg); err != nil {
			_, err := conn.Write([]byte("JSON Error: " + err.Error() + "\n"))
			if err != nil {
				return nil, err
			}
			continue
		}

		if err := validate.Struct(msg); err != nil {
			_, err := conn.Write([]byte("Validation Error: " + err.Error() + "\n"))
			if err != nil {
				return nil, err
			}
			continue
		}

		userId, err := c.parseJWTToken(ctx, msg.UserCredential)
		if err != nil {
			_, err := conn.Write([]byte("Token Error: " + err.Error() + "\n"))
			if err != nil {
				return nil, err
			}
			continue
		}

		if err := c.repo.Set(ctx, userId.String(), conn.RemoteAddr().String()); err != nil {
			_, err := conn.Write([]byte("Repository Error: " + err.Error() + "\n"))
			if err != nil {
				return nil, err
			}
			continue
		}

		return &msg, nil
	}

	if err := scanner.Err(); err != nil {
		_, err := conn.Write([]byte("Connection Error: " + err.Error() + "\n"))
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &msg, nil
}

func (c *MessageController) writeTCPRequest(ctx context.Context, message *entity.TCPRequest) error {
	colleagueAddr, err := c.repo.Get(ctx, message.UserCredential.ColleagueId.String())
	if err != nil {
		return err
	}

	colleagueConn, err := net.Dial(c.cfg.Server.Type, colleagueAddr.(string))
	if err != nil {
		return err
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = colleagueConn.Write(messageJSON)
	if err != nil {
		return err
	}

	return nil
}

func (c *MessageController) parseJWTToken(ctx context.Context, user entity.UserCredential) (uuid.UUID, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(user.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.cfg.AuthenticationConfiguration.AccessSecretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userIdStr, ok := claims[c.cfg.Claims.KeyForId].(string)
	if !ok {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}
