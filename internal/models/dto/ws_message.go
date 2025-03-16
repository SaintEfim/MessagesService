package dto

import (
	"time"

	"github.com/google/uuid"
)

type WsMessages struct {
	SenderId   uuid.UUID `json:"senderId" validate:"required"`
	ReceiverId uuid.UUID `json:"receiverId" validate:"required"`
	Text       string    `json:"text"`
	CreateAt   time.Time `json:"createAt"`
}
