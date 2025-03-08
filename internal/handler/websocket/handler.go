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

	request := &dto.SendMessage{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.controller.SendMessage(ctx, request); err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.handleError(nil, "Upgrader error", err)
		return
	}

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
	if err := h.controller.Connect(ctx, clientModel, transfer); err != nil {
		h.handleError(conn, "Invalid connect", err)
		return
	}

	_ = conn.WriteJSON(&dto.ResponseMessage{Text: "Success connect!"})
}

func (h *Handler) handleError(conn *websocket.Conn, message string, err error) {
	h.logger.Error(message, zap.Error(err))
	if conn != nil {
		_ = conn.WriteJSON(&dto.ResponseMessage{Error: message + ": " + err.Error()})
	}
}
