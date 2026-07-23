package kafka

import "time"

const (
	EventCourseCreated   = "course.created"
	EventCourseUpdated   = "course.updated"
	EventCourseDeleted   = "course.deleted"
	EventLessonCreated   = "lesson.created"
	EventLessonPublished = "lesson.published"
	EventStudentEnrolled = "student.enrolled"
)

type CourseCreatedEvent struct {
	Event     string    `json:"event"`
	CourseID  uint      `json:"course_id"`
	TeacherID uint      `json:"teacher_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type CourseUpdatedEvent struct {
	Event       string    `json:"event"`
	CourseID    uint      `json:"course_id"`
	TeacherID   uint      `json:"teacher_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CourseDeletedEvent struct {
	Event     string    `json:"event"`
	CourseID  uint      `json:"course_id"`
	TeacherID uint      `json:"teacher_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LessonCreatedEvent struct {
	Event       string    `json:"event"`
	LessonID    uint      `json:"lesson_id"`
	CourseID    uint      `json:"course_id"`
	Title       string    `json:"title"`
	OrderNumber int       `json:"order_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type LessonPublishedEvent struct {
	Event       string    `json:"event"`
	LessonID    uint      `json:"lesson_id"`
	CourseID    uint      `json:"course_id"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
}

type StudentEnrolledEvent struct {
	Event     string    `json:"event"`
	CourseID  uint      `json:"course_id"`
	StudentID uint      `json:"student_id"`
	CreatedAt time.Time `json:"created_at"`
}
