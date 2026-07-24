package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/gateway/middleware"
	"github.com/joho/godotenv"
)

func reverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("url parse failed: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}

}

func main() {
	r := gin.Default()

	if val := os.Getenv("IS_DOCKER"); val == "" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal("failed to parse envs")
		}
	}

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("user service url is not specified")
	}

	r.POST("/auth/login", reverseProxy(userServiceURL))
	r.POST("/auth/register", reverseProxy(userServiceURL))
	r.GET("/auth/me",
		middleware.AuthMiddleware,
		middleware.InjectHeaders(),
		reverseProxy(userServiceURL),
	)

	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware, middleware.InjectHeaders())

	protected.Any("/users/*any",
		reverseProxy(userServiceURL),
	)

	if err := r.Run(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
