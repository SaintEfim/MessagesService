package dto

import "github.com/google/uuid"

type CreateAction struct {
	Id uuid.UUID `json:"id"`
}
