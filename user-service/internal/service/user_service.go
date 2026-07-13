package service

import (
	// "github.com/itsvagapov/team-LMS/user-service/internal/model"
	"github.com/itsvagapov/team-LMS/user-service/internal/repository"
)

type UserService interface {
	ChangeRole()
}

type userService struct {
	users repository.UserRepository
}
