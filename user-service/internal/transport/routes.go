package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/user-service/internal/service"
)

func RegisterRouts(r *gin.Engine, auth service.AuthService, user service.UserService) {
	authHandler := NewAuthHandler(auth)
	userHandler := NewUserHandler(user, auth)

	authHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
}
