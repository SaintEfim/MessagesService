package websocket

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
	r.HandleFunc("/api/v1/message/connect", h.Connect).Methods("GET")
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.handleError(nil, "Upgrader error", err)
		return
	}
	defer conn.Close()

	h.handleMessages(ctx, conn)
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.handleError(nil, "Upgrader error", err)
		return
	}
	defer conn.Close()

	_, req, err := conn.ReadMessage()
	if err != nil {
		h.handleError(conn, "Failed to read message", err)
		return
	}

	clientModel := &dto.ConnectClient{}
	if err := json.Unmarshal(req, clientModel); err != nil {
		h.handleError(conn, "Invalid JSON format", err)
		return
	}

	transfer := websocketTransfer.NewWebSocketConnection(conn)
	h.controller.Connect(ctx, clientModel, transfer)
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
			h.handleError(conn, "Error sending message", err)
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
	if err := json.Unmarshal(msg, request); err != nil {
		return nil, err
	}

	return request, nil
}

func (h *Handler) handleError(conn *websocket.Conn, message string, err error) {
	h.logger.Error(message, zap.Error(err))
	if conn != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(message+": "+err.Error()))
	}
}
