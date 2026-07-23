package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		role := c.GetHeader("X-User-Role")
		if userID == "" || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
			return
		}

		parsedID, err := strconv.ParseUint(userID, 10, 64)
		if err != nil || parsedID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}

		c.Set("user_id", uint(parsedID))
		c.Set("user_role", role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		roleValue, _ := role.(string)
		if _, ok := allowed[roleValue]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
}
