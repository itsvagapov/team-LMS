package repository

import (
	"errors"

	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id uint) (*model.User, error)
	ChangeRole(id uint, role model.UserRole) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (r *gormUserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User

	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *gormUserRepository) ChangeRole(id uint, role model.UserRole) error {
	return r.db.
		Model(&model.User{}).
		Where("id = ?", id).
		Update("role", role).
		Error
}