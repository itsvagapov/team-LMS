package repository

import (
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"gorm.io/gorm"
)

type CourseEventRepository interface {
	Create(event *models.CourseEvent) error
}

type courseEventRepository struct {
	db *gorm.DB
}

func NewCourseEventRepository(
	db *gorm.DB,
) CourseEventRepository {
	return &courseEventRepository{db: db}
}

func (r *courseEventRepository) Create(
	event *models.CourseEvent,
) error {
	return r.db.Create(event).Error
}
