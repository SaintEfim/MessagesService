package dto

import "time"

type Messages struct {
	Text     string    `json:"text"`
	CreateAt time.Time `json:"create_at"`
}
