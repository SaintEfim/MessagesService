package dto

import "github.com/google/uuid"

type SendMessage struct {
	ChatId     uuid.UUID `json:"chat_id"`
	SenderId   uuid.UUID `json:"sender_id" validate:"required"`
	ReceiverId uuid.UUID `json:"receiver_id" validate:"required"`
	Text       string    `json:"text" validate:"gt=0,lte=100"`
}
