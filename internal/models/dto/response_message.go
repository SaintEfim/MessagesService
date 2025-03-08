package dto

import "time"

type ResponseMessage struct {
	Text      string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Error     string    `json:"error"`
}
