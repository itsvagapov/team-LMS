package repository

import (
	"activity-analytics-service/internal/model"

	"gorm.io/gorm"
)

type ActivityEventRepository struct {
	db *gorm.DB
}

func NewActivityEventRepository(db *gorm.DB) *ActivityEventRepository {
	return &ActivityEventRepository{db: db}
}

func (r *ActivityEventRepository) Create(event *model.ActivityEvent) error {
	return r.db.Create(event).Error
}

func (r *ActivityEventRepository) List(limit, offset int) ([]model.ActivityEvent, error) {
	var events []model.ActivityEvent
	err := r.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&events).Error
	return events, err
}

func (r *ActivityEventRepository) ListByUser(userID uint, limit, offset int) ([]model.ActivityEvent, error) {
	var events []model.ActivityEvent
	err := r.db.Where("actor_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

func (r *ActivityEventRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.ActivityEvent{}).Count(&count).Error
	return count, err
}
