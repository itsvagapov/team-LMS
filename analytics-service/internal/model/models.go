package model

import "time"

type ActivityEvent struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EventType     string    `gorm:"not null;index" json:"event_type"`
	SourceService string    `gorm:"not null;index" json:"source_service"`
	ActorID       *uint     `gorm:"index" json:"actor_id,omitempty"`
	EntityType    string    `gorm:"index" json:"entity_type"`
	EntityID      *uint     `gorm:"index" json:"entity_id,omitempty"`
	Payload       string    `gorm:"type:jsonb;not null" json:"payload"`
	CreatedAt     time.Time `gorm:"not null;index" json:"created_at"`
}

type CourseStats struct {
	ID                       uint      `gorm:"primaryKey" json:"id"`
	CourseID                 uint      `gorm:"uniqueIndex;not null" json:"course_id"`
	StudentsCount            int       `gorm:"not null;default:0" json:"students_count"`
	LessonsCount             int       `gorm:"not null;default:0" json:"lessons_count"`
	HomeworksCount           int       `gorm:"not null;default:0" json:"homeworks_count"`
	SubmissionsCount         int       `gorm:"not null;default:0" json:"submissions_count"`
	ReviewedSubmissionsCount int       `gorm:"not null;default:0" json:"reviewed_submissions_count"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type UserActivityStats struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	RegisteredAt         *time.Time `json:"registered_at,omitempty"`
	CoursesEnrolledCount int        `gorm:"not null;default:0" json:"courses_enrolled_count"`
	SubmissionsCount     int        `gorm:"not null;default:0" json:"submissions_count"`
	LastActivityAt       *time.Time `json:"last_activity_at,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type Dashboard struct {
	EventsCount              int64 `json:"events_count"`
	UsersCount               int64 `json:"users_count"`
	CoursesCount             int64 `json:"courses_count"`
	StudentsEnrolledCount    int64 `json:"students_enrolled_count"`
	SubmissionsCount         int64 `json:"submissions_count"`
	ReviewedSubmissionsCount int64 `json:"reviewed_submissions_count"`
}
