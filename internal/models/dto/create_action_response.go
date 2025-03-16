package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateActionResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
