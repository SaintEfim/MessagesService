package dto

import "time"

type WsMessages struct {
	Text     string    `json:"text"`
	CreateAt time.Time `json:"create_at"`
}
