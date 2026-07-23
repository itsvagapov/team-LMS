package repository

import (
	"time"

	"activity-analytics-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserActivityStatsRepository struct {
	db *gorm.DB
}

func NewUserActivityStatsRepository(db *gorm.DB) *UserActivityStatsRepository {
	return &UserActivityStatsRepository{db: db}
}

func (r *UserActivityStatsRepository) GetByUserID(userID uint) (*model.UserActivityStats, error) {
	var stats model.UserActivityStats
	err := r.db.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *UserActivityStatsRepository) Ensure(userID uint) error {
	stats := model.UserActivityStats{UserID: userID}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&stats).Error
}

func (r *UserActivityStatsRepository) MarkRegistered(userID uint, registeredAt time.Time) error {
	if err := r.Ensure(userID); err != nil {
		return err
	}
	return r.db.Model(&model.UserActivityStats{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"registered_at":    registeredAt,
			"last_activity_at": registeredAt,
		}).Error
}

func (r *UserActivityStatsRepository) Increment(userID uint, column string, value int, lastActivityAt time.Time) error {
	if err := r.Ensure(userID); err != nil {
		return err
	}
	return r.db.Model(&model.UserActivityStats{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			column:             gorm.Expr(column+" + ?", value),
			"last_activity_at": lastActivityAt,
		}).Error
}

func (r *UserActivityStatsRepository) Touch(userID uint, lastActivityAt time.Time) error {
	if err := r.Ensure(userID); err != nil {
		return err
	}
	return r.db.Model(&model.UserActivityStats{}).
		Where("user_id = ?", userID).
		Update("last_activity_at", lastActivityAt).Error
}

func (r *UserActivityStatsRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.UserActivityStats{}).Count(&count).Error
	return count, err
}
