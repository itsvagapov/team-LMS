package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/course-service/internal/config"
	"github.com/itsvagapov/team-LMS/course-service/internal/handler"
	coursekafka "github.com/itsvagapov/team-LMS/course-service/internal/kafka"
	"github.com/itsvagapov/team-LMS/course-service/internal/repository"
	"github.com/itsvagapov/team-LMS/course-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	log.Println("database connected successfully")

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	brokers := strings.Split(kafkaBrokers, ",")
	kafkaProducer := coursekafka.NewProducer(brokers)

	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("failed to close kafka producer: %v", err)
		}
	}()

	courseRepository := repository.NewCourseRepository(config.DB)
	courseService := service.NewCourseService(courseRepository, kafkaProducer)
	courseHandler := handler.NewCourseHandler(courseService)

	lessonRepository := repository.NewLessonRepository(config.DB)
	lessonService := service.NewLessonService(courseService, lessonRepository, kafkaProducer)
	lessonHandler := handler.NewLessonHandler(lessonService)

	courseStudentRepository := repository.NewCourseStudentRepository(config.DB)

	courseEventRepository := repository.NewCourseEventRepository(config.DB)

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "course-service"
	}

	kafkaConsumer := coursekafka.NewConsumer(
		brokers,
		kafkaGroupID,
		courseEventRepository,
	)

	consumerCtx, cancelConsumer := context.WithCancel(
		context.Background(),
	)

	go kafkaConsumer.Run(consumerCtx)

	defer func() {
		cancelConsumer()

		if err := kafkaConsumer.Close(); err != nil {
			log.Printf("failed to close kafka consumer: %v", err)
		}
	}()

	enrollmentService := service.NewEnrollmentService(courseService, courseStudentRepository, kafkaProducer)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService)

	internalService := service.NewInternalService(courseService, lessonService, courseStudentRepository)
	internalHandler := handler.NewInternalHandler(internalService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "course-service",
		})
	})

	courseHandler.RegisterRoutes(router)
	lessonHandler.RegisterRoutes(router)
	enrollmentHandler.RegisterRoutes(router)
	internalHandler.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("starting course-service on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("failed to start course-service:", err)
	}
}
