package websocket

import (
	"MessagesService/internal/models/interfaces"

	"github.com/gorilla/websocket"
)

type WebSocketTransfer struct {
	conn *websocket.Conn
}

func NewWebSocketConnection(conn *websocket.Conn) interfaces.Transfer {
	return &WebSocketTransfer{conn: conn}
}

func (ws *WebSocketTransfer) TransferData(data interface{}) error {
	return ws.conn.WriteJSON(data)
}
func (ws *WebSocketTransfer) TransferDataText(text string) error {
	return ws.conn.WriteMessage(websocket.TextMessage, []byte(text))
}
