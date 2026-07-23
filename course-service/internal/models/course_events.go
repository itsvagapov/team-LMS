package models

import "time"

type CourseEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventType string    `gorm:"not null;index" json:"event_type"`
	Payload   string    `gorm:"type:jsonb;not null" json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}
