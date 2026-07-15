package repository

import (
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id uint) (*model.UserResponse, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) GetByID(id uint) (*model.UserResponse, error) {
	return nil, nil
}