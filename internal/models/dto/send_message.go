package dto

import "github.com/google/uuid"

type SendMessage struct {
	SenderID   uuid.UUID `json:"sender_id" validate:"required"`
	ReceiverID uuid.UUID `json:"receiver_id" validate:"required"`
	Message    string    `json:"message" validate:"gt=0,lte=4096"`
}
