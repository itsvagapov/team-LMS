package models

import "time"

type Lesson struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseID    uint      `gorm:"not null;index;uniqueIndex:idx_course_order" json:"course_id"`
	Title       string    `gorm:"type:varchar(50);not null" json:"title"`
	Content     string    `gorm:"type:varchar(2000);not null" json:"content"`
	OrderNumber int       `gorm:"not null;uniqueIndex:idx_course_order" json:"order_number"`
	IsPublished bool      `gorm:"not null;default:false" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
