package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/user-service/internal/middleware"
)

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	protected := r.Group("")
	protected.Use(middleware.GatewayHeadersMiddleware())

	unprotected := r.Group("")

	authNoDefense := unprotected.Group("/auth")
	{
		authNoDefense.POST("/register", h.Register)
		authNoDefense.POST("/login", h.Login)
	}

	authDefense := protected.Group("/auth")
	{
		authDefense.GET("/me", h.Me)
	}
}
