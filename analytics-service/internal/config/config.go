package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort       string
	DatabaseDSN    string
	KafkaBrokers   []string
	KafkaGroupID   string
	KafkaTopics    []string
	AnalyticsTopic string
}

func Load() Config {
	godotenv.Load(".env")

	return Config{
		HTTPPort:       getEnv("HTTP_PORT", "8084"),
		DatabaseDSN:    getEnv("DATABASE_DSN", "host=postgres user=postgres password=postgres dbname=analytics_db port=5432 sslmode=disable"),
		KafkaBrokers:   splitEnv("KAFKA_BROKERS", "kafka:9092"),
		KafkaGroupID:   getEnv("KAFKA_GROUP_ID", "activity-analytics-service"),
		KafkaTopics:    splitEnv("KAFKA_TOPICS", "users.events,courses.events,homework.events"),
		AnalyticsTopic: getEnv("ANALYTICS_TOPIC", "analytics.events"),
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitEnv(key, fallback string) []string {
	raw := getEnv(key, fallback)
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}
