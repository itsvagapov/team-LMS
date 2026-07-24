package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/gateway/internal/jwt"
)

func AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing authorization header",
		})
		ctx.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		ctx.Abort()
		return
	}

	ctx.Set("userID", claims.UserID)
	ctx.Set("userRole", claims.Role)

	ctx.Next()
}

func InjectHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			log.Fatal("userID missing in context")
			ctx.Abort()
			return
		}

		role, exists := ctx.Get("userRole")
		if !exists {
			log.Fatal("userRole missing in context")
			ctx.Abort()
			return
		}

		ctx.Request.Header.Set(
			"X-User-ID",
			fmt.Sprintf("%v", userID),
		)

		ctx.Request.Header.Set(
			"X-User-Role",
			fmt.Sprintf("%v", role),
		)

		ctx.Next()
	}
}
