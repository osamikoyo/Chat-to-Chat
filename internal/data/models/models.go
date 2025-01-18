package models

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
