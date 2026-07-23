package repository

import (
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(lesson *models.Lesson) error
	FindByID(id uint) (*models.Lesson, error)
	FindByCourseID(courseID uint, publishedOnly bool) ([]models.Lesson, error)
	GetNextOrderNumber(courseID uint) (int, error)
	Update(lesson *models.Lesson) error
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) Create(lesson *models.Lesson) error {
	return r.db.Create(lesson).Error
}

func (r *lessonRepository) GetNextOrderNumber(courseID uint) (int, error) {
	var nextOrderNumber int

	err := r.db.
		Model(&models.Lesson{}).
		Where("course_id = ?", courseID).
		Select("COALESCE(MAX(order_number), 0) + 1").
		Scan(&nextOrderNumber).
		Error

	if err != nil {
		return 0, err
	}
	return nextOrderNumber, nil
}

func (r *lessonRepository) FindByCourseID(courseID uint, publishedOnly bool) ([]models.Lesson, error) {
	var lessons []models.Lesson

	query := r.db.
		Where("course_id = ?", courseID).
		Order("order_number ASC")

	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}

	if err := query.Find(&lessons).Error; err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *lessonRepository) FindByID(id uint) (*models.Lesson, error) {
	var lesson models.Lesson

	if err := r.db.First(&lesson, id).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *lessonRepository) Update(lesson *models.Lesson) error {
	return r.db.Save(lesson).Error
}
