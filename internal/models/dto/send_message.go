package dto

import "github.com/google/uuid"

type SendMessage struct {
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Message    string    `json:"message"`
}
