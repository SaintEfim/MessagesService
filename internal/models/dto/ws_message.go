package dto

import (
	"time"

	"github.com/google/uuid"
)

type WsMessages struct {
	ChatId     uuid.UUID `json:"chat_id"`
	SenderId   uuid.UUID `json:"senderId"`
	ReceiverId uuid.UUID `json:"receiverId"`
	Text       string    `json:"text"`
	CreateAt   time.Time `json:"createAt"`
}
