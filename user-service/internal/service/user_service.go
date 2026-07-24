package service

import (
	"github.com/itsvagapov/team-LMS/user-service/internal/kafkabro"
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"github.com/itsvagapov/team-LMS/user-service/internal/repository"
)

type UserService interface {
	ChangeRole(id uint, role model.UserRole) error
}

type userService struct {
	users    repository.UserRepository
	auth     repository.AuthRepository
	producer kafkabro.Producer
}


func NewUserService(auth repository.AuthRepository, users repository.UserRepository, producer kafkabro.Producer) UserService {
	return &userService{
		auth:     auth,
		users:    users,
		producer: producer,
	}
}

func (s *userService) ChangeRole(id uint, role model.UserRole) error {
	user, err := s.users.GetByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	if user.Role == role {
		return ErrRoleAlreadyAssigned
	}

	return s.users.ChangeRole(id, role)
}