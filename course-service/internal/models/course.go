package models

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(50);not null" json:"title"`
	Description string    `gorm:"type:varchar(200)" json:"description"`
	TeacherID   uint      `gorm:"not null;index" json:"teacher_id"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
