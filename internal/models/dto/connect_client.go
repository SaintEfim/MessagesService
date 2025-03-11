package dto

import "github.com/google/uuid"

type ConnectClient struct {
	Id uuid.UUID `json:"Id" validate:"required"`
}
