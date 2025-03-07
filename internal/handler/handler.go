package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"MessagesService/config"
	websocketTransfer "MessagesService/internal/delivery/websocket"
	"MessagesService/internal/models/dto"
	"MessagesService/internal/models/interfaces"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	controller interfaces.Controller
	logger     *zap.Logger
	upgrader   *websocket.Upgrader
	cfg        *config.Config
}

func NewHandler(
	controller interfaces.Controller,
	logger *zap.Logger,
	upgrader *websocket.Upgrader,
	cfg *config.Config) interfaces.Handler {

	return &Handler{
		controller: controller,
		logger:     logger,
		upgrader:   upgrader,
		cfg:        cfg,
	}
}

func (h *Handler) ConfigureRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/message", h.SendMessage).Methods("POST")
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("upgrader error", zap.Error(err))
		return
	}
	defer conn.Close()

	h.handleMessages(ctx, conn)
}

func (h *Handler) handleMessages(ctx context.Context, conn *websocket.Conn) {
	for {
		request, err := h.readMessage(conn)
		if err != nil {
			h.handleError(conn, "Message read error", err)
			break
		}

		transfer := websocketTransfer.NewWebSocketConnection(conn)

		if err := h.controller.SendMessage(ctx, request, transfer); err != nil {
			h.handleError(conn, "Error send message:", err)
			break
		}
	}
}

func (h *Handler) readMessage(conn *websocket.Conn) (*dto.SendMessage, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	request := &dto.SendMessage{}
	if err := json.Unmarshal(msg, &request); err != nil {
		return nil, err
	}

	return request, nil
}

func (h *Handler) handleError(conn *websocket.Conn, message string, err error) {
	h.logger.Error(message, zap.Error(err))
	_ = conn.WriteMessage(websocket.TextMessage, []byte(message+": "+err.Error()))
}
