package config

import (
	"fmt"
	"os"

	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbName,
		port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	if err := DB.AutoMigrate(
		&models.Course{},
		&models.Lesson{},
		&models.CourseStudent{},
		&models.CourseEvent{},
	); err != nil {
		return err
	}

	return nil
}
