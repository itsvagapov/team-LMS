package main

import (
	"log"
	"os"
	"time"

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

	db, err := config.NewDBConn(cfg.DB)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("failed to run model migrations: %v", err)
	}

	err = kafkabro.CreateTopic(cfg.Kafka, kafkabro.UsersEventsTopic)

	for i := 0; i < 30; i++ {
		err = kafkabro.CreateTopic(cfg.Kafka, kafkabro.UsersEventsTopic)
		if err == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}

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

	err = authService.CreateSuperAdmin("admin", "admin@mail.ru", "123456")
	if err != nil {
		log.Fatal("failed to create super-admin")
	}

	if err := r.Run(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("failed to run the HTTP server")
	}
}
