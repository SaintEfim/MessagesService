package dto

import "github.com/google/uuid"

type Connect struct {
	Id uuid.UUID `json:"Id" validate:"required"`
}
