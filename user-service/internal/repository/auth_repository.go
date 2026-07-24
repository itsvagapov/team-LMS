package repository

import (
	"errors"

	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *model.User) error
	ExistsByEmail(email string) (bool, error)
	GetByEmail(email string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
}

type gormAuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &gormAuthRepository{db: db}
}

func (r *gormAuthRepository) CreateUser(user *model.User) error {
	if user == nil {
		return ErrUserIsNil
	}

	return r.db.Create(user).Error
}

func (r *gormAuthRepository) ExistsByEmail(email string) (bool, error) {
	var count int64

	err := r.db.Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *gormAuthRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}


func (r *gormAuthRepository) GetByID(id uint) (*model.User, error) {
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
