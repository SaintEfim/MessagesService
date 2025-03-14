package websocket

import (
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
	r.HandleFunc("/api/v1/message/connect", h.Connect)
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := &dto.SendMessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		h.logger.Error("Error decoding request", zap.Error(err))
		return
	}

	res, err := h.controller.SendMessage(ctx, request)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		h.logger.Error("Error sending message", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.Error("Error encoding response", zap.Error(err))
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.handleError(nil, "Upgrader error", err)
		h.logger.Error("Error upgrading connection", zap.Error(err))
		return
	}

	_, req, err := conn.ReadMessage()
	if err != nil {
		h.handleError(conn, "Failed to read message", err)
		h.logger.Error("Error reading message", zap.Error(err))
		return
	}

	clientModel := &dto.ConnectClientRequest{}
	if err := json.Unmarshal(req, clientModel); err != nil {
		h.handleError(conn, "Invalid JSON format", err)
		h.logger.Error("Error unmarshalling message", zap.Error(err))
		return
	}

	transfer := websocketTransfer.NewWebSocketConnection(conn)
	if err := h.controller.Connect(ctx, clientModel, transfer); err != nil {
		h.handleError(conn, "Invalid connect", err)
		h.logger.Error("Error connecting to client", zap.Error(err))
		return
	}

	_ = conn.WriteMessage(websocket.TextMessage, []byte("Success connect!"))
}

func (h *Handler) handleError(conn *websocket.Conn, message string, err error) {
	h.logger.Error(message, zap.Error(err))
	if conn != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(message+": "+err.Error()))
		h.logger.Warn("Error sending message", zap.Error(err))
	}
}
