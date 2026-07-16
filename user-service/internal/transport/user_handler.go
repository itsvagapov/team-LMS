package transport

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/user-service/internal/middleware"
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"github.com/itsvagapov/team-LMS/user-service/internal/service"
)

type UserHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	protected := r.Group("")
	protected.Use(middleware.GatewayHeadersMiddleware())
	
	userDefense := protected.Group("/users")
	{
		userDefense.PATCH("/:id/role", h.ChangeRole)
	}
}

func (h *UserHandler) ChangeRole(c *gin.Context) {
	userRole, exists := c.Get("userRole")
	if !exists {
		log.Println("missing user role in context")
		c.Status(http.StatusInternalServerError)
		return
	}

	if userRole != model.RoleAdmin {
		log.Println("insufficient privileges")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you do not have sufficient permissions",
		})
		return
	}

	userIDStr := c.Param("id")
	if userIDStr == "" {
		log.Println("user id is not specified")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id is not specified",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("user id cast failed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var req model.ChangeRoleRequest

	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Println("invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.authService.GetUserByID(uint(userID))
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if user.Role == req.Role {
		log.Println("user is already have such role")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("user is already: %v", req.Role),
		})
		return 
	}

	err = h.userService.ChangeRole(uint(userID), req.Role)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}