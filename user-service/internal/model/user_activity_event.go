package model

import "time"

type UserActivityEvent struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint `gorm:"not null;index" json:"user_id"`

	EventType UserActivityEventType `gorm:"type:varchar(100);not null" json:"event_type"`

	SourceService SourceService `gorm:"type:varchar(100);not null" json:"source_service"`

	Payload string `gorm:"type:text" json:"payload"`

	CreatedAt time.Time `json:"created_at"`
}