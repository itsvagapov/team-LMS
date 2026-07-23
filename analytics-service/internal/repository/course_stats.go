package repository

import (
	"activity-analytics-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourseStatsRepository struct {
	db *gorm.DB
}

func NewCourseStatsRepository(db *gorm.DB) *CourseStatsRepository {
	return &CourseStatsRepository{db: db}
}

func (r *CourseStatsRepository) GetByCourseID(courseID uint) (*model.CourseStats, error) {
	var stats model.CourseStats
	err := r.db.Where("course_id = ?", courseID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *CourseStatsRepository) Ensure(courseID uint) error {
	stats := model.CourseStats{CourseID: courseID}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&stats).Error
}

func (r *CourseStatsRepository) Increment(courseID uint, column string, value int) error {
	if err := r.Ensure(courseID); err != nil {
		return err
	}
	return r.db.Model(&model.CourseStats{}).
		Where("course_id = ?", courseID).
		Updates(map[string]any{
			column: columnExpr(column, value),
		}).
		Error
}

func (r *CourseStatsRepository) DashboardSums() (students, submissions, reviewed int64, err error) {
	var result struct {
		Students    int64
		Submissions int64
		Reviewed    int64
	}
	err = r.db.Model(&model.CourseStats{}).
		Select("coalesce(sum(students_count), 0) as students, coalesce(sum(submissions_count), 0) as submissions, coalesce(sum(reviewed_submissions_count), 0) as reviewed").
		Scan(&result).Error
	return result.Students, result.Submissions, result.Reviewed, err
}

func (r *CourseStatsRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.CourseStats{}).Count(&count).Error
	return count, err
}

func columnExpr(column string, value int) any {
	return gorm.Expr(column+" + ?", value)
}
