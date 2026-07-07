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
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
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

	resp, err := h.userService.Register(c.Request.Context(), req)
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

	resp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user id from context",
		})
		return
	}

	userIDCasted, ok := userID.(int)
	if !ok {
		c.Status(http.StatusInternalServerError)
		log.Println("failed cast userID from any to int")
		return
	}

	user, err := h.userService.GetUserByID(userIDCasted)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, user)
}
