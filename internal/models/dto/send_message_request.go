package dto

import "github.com/google/uuid"

type SendMessageRequest struct {
	ChatId     uuid.UUID `json:"chatId"`
	SenderId   uuid.UUID `json:"senderId" validate:"required"`
	ReceiverId uuid.UUID `json:"receiverId" validate:"required"`
	Text       string    `json:"text" validate:"gt=0,lte=100"`
}
