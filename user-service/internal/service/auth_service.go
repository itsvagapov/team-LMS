package service

import (
	"context"
	"encoding/json"
	"log"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/itsvagapov/team-LMS/user-service/internal/jwt"
	"github.com/itsvagapov/team-LMS/user-service/internal/kafkabro"
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"github.com/itsvagapov/team-LMS/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req *model.RegisterRequest) (*model.UserResponse, error)
	LoginUser(req model.LoginRequest) (*model.AuthResponse, error)
	GetUserByID(id uint) (*model.UserResponse, error)
}

type authService struct {
	auth     repository.AuthRepository
	users    repository.UserRepository
	producer kafkabro.Producer
}

func NewAuthService(auth repository.AuthRepository, users repository.UserRepository, producer kafkabro.Producer) AuthService {
	return &authService{
		auth:     auth,
		users:    users,
		producer: producer,
	}
}

func (s *authService) RegisterUser(req *model.RegisterRequest) (*model.UserResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	exists, err := s.auth.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         model.RoleStudent,
	}

	if err := s.auth.CreateUser(&user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *authService) LoginUser(req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.auth.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidLoginOrPassword
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, ErrInvalidLoginOrPassword
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(kafkabro.UserLoggedInPayload{
		Role: user.Role,
	})
	if err != nil {
		return nil, err
	}

	event := kafkabro.UserActivityEventMessage{
		UserID:        user.ID,
		EventType:     kafkabro.EventUserLoggedIn,
		SourceService: kafkabro.ServiceUser,
		Payload:       string(payload),
		CreatedAt:     time.Now().UTC(),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	err = s.producer.Send(
		context.Background(),
		eventBytes,
	)
	if err != nil {
		log.Println("ОШИБКА КАФКИ: ",err)
		return nil, err
	}

	return &model.AuthResponse{Token: token}, nil
}

func (s *authService) GetUserByID(id uint) (*model.UserResponse, error) {
	user, err := s.auth.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func validateRegisterRequest(req *model.RegisterRequest) error {
	name := req.Name

	if req.Name == "" {
		return ErrNameRequired
	}

	if len(name) < 2 || len(name) > 30 {
		return ErrInvalidNameLength
	}

	if strings.ContainsRune(name, ' ') {
		return ErrNameContainsSpaces
	}

	for _, r := range name {
		if unicode.IsDigit(r) {
			return ErrNameContainDigits
		}
	}

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || r == '-' {
			continue
		}
		return ErrInvalidCharacters
	}

	if req.Email == "" {
		return ErrEmailRequired
	}

	if len(req.Email) > 254 {
		return ErrInvalidEmail
	}

	email := req.Email

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return ErrInvalidEmail
	}

	if strings.ContainsRune(email, ' ') {
		return ErrEmailContainsSpaces
	}

	email = strings.ToLower(email)

	if req.Password == "" {
		return ErrPasswordRequired
	}

	if strings.ContainsRune(req.Password, ' ') {
		return ErrPasswordContainsSpaces
	}

	return nil
}
