package dto

import "github.com/google/uuid"

type ConnectClientRequest struct {
	Id uuid.UUID `json:"Id" validate:"required"`
}
