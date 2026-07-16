package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/user-service/internal/config"
	"github.com/itsvagapov/team-LMS/user-service/internal/kafkabro"
	"github.com/itsvagapov/team-LMS/user-service/internal/model"
	"github.com/itsvagapov/team-LMS/user-service/internal/repository"
	"github.com/itsvagapov/team-LMS/user-service/internal/service"
	"github.com/itsvagapov/team-LMS/user-service/internal/transport"
)

func main() {
	r := gin.Default()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config failed to initialize: %v", err)
	}

	log.Println(cfg)

	db, err := config.NewDBConn(cfg.DB)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("failed to run model migrations: %v", err)
	}

	err = kafkabro.CreateTopic(cfg.Kafka, kafkabro.UsersEventsTopic)
	if err != nil {
		log.Fatalf("failed to create topic: %v", err)
	}

	producer := kafkabro.NewProducer(
		[]string{cfg.Kafka.Brokers[0]},
		kafkabro.UsersEventsTopic,
	)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, userRepo, *producer)
	userService := service.NewUserService(authRepo, userRepo, *producer)
	transport.RegisterRouts(r, authService, userService)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("failed to run the HTTP server")
	}
}
