package repository

import (
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"gorm.io/gorm"
)

type CourseStudentRepository interface {
	Create(courseStudent *models.CourseStudent) error
	Exists(courseID uint, studentID uint) (bool, error)
	FindByCourseID(courseID uint) ([]models.CourseStudent, error)
}

type courseStudentRepository struct {
	db *gorm.DB
}

func NewCourseStudentRepository(db *gorm.DB) CourseStudentRepository {
	return &courseStudentRepository{db: db}
}

func (r *courseStudentRepository) Create(courseStudent *models.CourseStudent) error {
	return r.db.Create(courseStudent).Error
}

func (r *courseStudentRepository) Exists(courseID uint, studentID uint) (bool, error) {
	var count int64

	if err := r.db.Model(&models.CourseStudent{}).
		Where("course_id = ? AND student_id = ?", courseID, studentID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *courseStudentRepository) FindByCourseID(courseID uint) ([]models.CourseStudent, error) {
	var courseStudents []models.CourseStudent

	if err := r.db.
		Where("course_id = ?", courseID).
		Order("created_at ASC").
		Find(&courseStudents).
		Error; err != nil {
		return nil, err
	}

	return courseStudents, nil
}
