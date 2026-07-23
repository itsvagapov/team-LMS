package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	UserIDContextKey   = "user_id"
	UserRoleContextKey = "user_role"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDHeader := ctx.GetHeader("X-User-ID")
		userRole := ctx.GetHeader("X-User-Role")

		userID, err := strconv.ParseUint(userIDHeader, 10, 64)
		if err != nil || userID == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing X-User-ID header",
			})
			return
		}

		if userRole == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-User-Role header",
			})
			return
		}

		roleAllowed := false

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		ctx.Set(UserIDContextKey, uint(userID))
		ctx.Set(UserRoleContextKey, userRole)

		ctx.Next()
	}
}
