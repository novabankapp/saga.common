package domain

import "time"

type ConsumedMessage struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewConsumedMessage(id string, createdAt time.Time) *ConsumedMessage {
	return &ConsumedMessage{
		id,
		createdAt,
	}
}
