package models

import "time"

type CourseStudent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null;index;uniqueIndex:idx_course_student" json:"course_id"`
	StudentID uint      `gorm:"not null;index;uniqueIndex:idx_course_student" json:"student_id"`
	CreatedAt time.Time `json:"created_at"`
}
