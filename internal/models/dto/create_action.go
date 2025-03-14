package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAction struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
