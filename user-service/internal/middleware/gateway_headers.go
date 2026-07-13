package middleware

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
)

func GatewayHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.GetHeader("X-User-ID")

		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrUserRoleNotFound,
			})
			return
		}

		role := c.GetHeader("X-User-Role")
		if !slices.Contains([]model.UserRole{model.RoleStudent, model.RoleAdmin, model.RoleTeacher}, model.UserRole(role)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": ErrUserRoleNotFound,
			})
			return
		}

		c.Set(ContextUserID, uint(id64))
		c.Set(ContextUserRole, model.UserRole(role))
		c.Next()
	}
}
