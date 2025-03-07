package dto

import "time"

type ReceiveMessage struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
