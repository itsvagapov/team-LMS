package repository

import (
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindAll() ([]models.Course, error)
	FindByID(id uint) (*models.Course, error)
	Update(course *models.Course) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) FindAll() ([]models.Course, error) {
	var courses []models.Course

	if err := r.db.
		Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) FindByID(id uint) (*models.Course, error) {
	var course models.Course

	if err := r.db.
		Where("id = ? AND is_active = ?", id, true).
		First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}
