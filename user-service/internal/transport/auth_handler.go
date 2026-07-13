package transport

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/user-service/internal/middleware"
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"github.com/itsvagapov/team-LMS/user-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.authService.RegisterUser(req)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.authService.LoginUser(req)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.Status(http.StatusInternalServerError)
		log.Println("user id missing in context")
		return
	}

	userIDCasted, ok := userID.(int)
	if !ok {
		c.Status(http.StatusInternalServerError)
		log.Println("failed cast userID from any to int")
		return
	}

	user, err := h.authService.GetUserByID(userIDCasted)
	if err != nil {
		
	}

	c.JSON(http.StatusOK, user)
}
