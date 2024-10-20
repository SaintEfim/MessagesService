package entity

import "github.com/google/uuid"

type UserCredential struct {
	Token       string    `json:"token" validate:"required,jwt"`
	ColleagueId uuid.UUID `json:"colleague_id"`
}
