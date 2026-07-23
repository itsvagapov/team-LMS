package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"activity-analytics-service/internal/config"
	"activity-analytics-service/internal/db"
	"activity-analytics-service/internal/handler"
	analyticskafka "activity-analytics-service/internal/kafka"
	"activity-analytics-service/internal/middleware"
	"activity-analytics-service/internal/model"
	"activity-analytics-service/internal/repository"
	"activity-analytics-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.AutoMigrate(
		&model.ActivityEvent{},
		&model.CourseStats{},
		&model.UserActivityStats{},
	); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	eventRepo := repository.NewActivityEventRepository(database)
	courseStatsRepo := repository.NewCourseStatsRepository(database)
	userStatsRepo := repository.NewUserActivityStatsRepository(database)

	producer := analyticskafka.NewProducer(cfg.KafkaBrokers, cfg.AnalyticsTopic)
	analyticsService := service.NewAnalyticsService(eventRepo, courseStatsRepo, userStatsRepo, producer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := analyticskafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID, cfg.KafkaTopics, analyticsService)
	go consumer.Start(ctx)

	router := gin.Default()
	apiHandler := handler.NewAnalyticsHandler(analyticsService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	protected := router.Group("/")
	protected.Use(middleware.UserContext())
	protected.GET("/activity/events", middleware.RequireRole("admin"), apiHandler.ListEvents)
	protected.GET("/activity/users/:id", middleware.RequireRole("admin"), apiHandler.GetUserActivity)
	protected.GET("/analytics/courses/:id", middleware.RequireRole("teacher", "admin"), apiHandler.GetCourseStats)
	protected.GET("/analytics/dashboard", middleware.RequireRole("admin"), apiHandler.GetDashboard)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("activity analytics service started on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down activity analytics service")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown http server: %v", err)
	}
	if err := producer.Close(); err != nil {
		log.Printf("close kafka producer: %v", err)
	}
	if err := consumer.Close(); err != nil {
		log.Printf("close kafka consumer: %v", err)
	}
}
