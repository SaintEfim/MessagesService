package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"

	"MessagesService/config"
	"MessagesService/internal/models/dto"
	"MessagesService/internal/models/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
	cfg    *config.Config
	repo   interfaces.CacheRepository
}

func NewController(logger *zap.Logger, cfg *config.Config, repo interfaces.CacheRepository) interfaces.Controller {
	return &Controller{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

func (c *Controller) MessageHandleRequest(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(conn)

	for {
		msg, err := c.readTCPRequest(ctx, scanner, conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.logger.Info("Client disconnected")
				break
			}

			_, err := conn.Write([]byte("Read message Error: " + err.Error() + "\n"))
			if err != nil {
				c.logger.Error("Error writing message " + err.Error())
				break
			}

			continue
		}

		switch msg.Operation {
		case dto.OperationInit:
			continue
		case dto.OperationSendMessage:
			if err := c.writeTCPRequest(ctx, msg); err != nil {
				c.logger.Error("Error writing message " + err.Error())
			}
		default:
			_, err := conn.Write([]byte("Operation not found\n"))
			if err != nil {
				c.logger.Error("Error writing message " + err.Error())
			}
		}
	}

	return nil
}

func (c *Controller) readTCPRequest(ctx context.Context, scanner *bufio.Scanner, conn net.Conn) (*dto.TCPRequest, error) {
	msg := dto.TCPRequest{}
	validate := validator.New()

	for scanner.Scan() {
		clientMessage := scanner.Text()

		if err := json.Unmarshal([]byte(clientMessage), &msg); err != nil {
			return nil, err
		}

		if err := validate.Struct(msg); err != nil {
			return nil, err
		}

		userId, err := c.parseJWTToken(ctx, msg.UserCredential)
		if err != nil {
			return nil, err
		}

		if err := c.repo.Set(ctx, userId.String(), conn.RemoteAddr().String()); err != nil {
			return nil, err
		}

		return &msg, nil
	}

	if err := scanner.Err(); err != nil {
		c.logger.Error("Error scanning TCP Request: " + err.Error())
		return nil, err
	}

	return nil, io.EOF
}

func (c *Controller) writeTCPRequest(ctx context.Context, message *dto.TCPRequest) error {
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

func (c *Controller) parseJWTToken(ctx context.Context, user dto.UserCredential) (uuid.UUID, error) {
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
