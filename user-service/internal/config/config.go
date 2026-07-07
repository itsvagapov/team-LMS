package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	Kafka  KafkaConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type KafkaConfig struct {
	Brokers []string
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, fmt.Errorf("failed to load env file: %w", err)
	}

	serverCfg := ServerConfig{
		Port: os.Getenv("SERVER_PORT"),
	}

	if serverCfg.Port == "" {
		return nil, errors.New("server port is unreachable")
	}

	dbCfg := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if dbCfg.Host == "" {
		return nil, errors.New("db host is not specified")
	}

	if dbCfg.Port == "" {
		return nil, errors.New("db port is not specified")
	}
	if dbCfg.User == "" {
		return nil, errors.New("db user is not specified")
	}
	if dbCfg.Password == "" {
		return nil, errors.New("db password is not specified")
	}
	if dbCfg.Name == "" {
		return nil, errors.New("db port is not specified")
	}
	if dbCfg.SSLMode == "" {
		return nil, errors.New("db sslmode is not specified")
	}

	ttl, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return nil, errors.New("jwt ttl invalid")
	}
	

	jwtCfg := JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		TTL: ttl,
	}

	if jwtCfg.Secret == "" {
		return nil, errors.New("jwt secret is not specified")
	}

	kafkaCfg := KafkaConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
	}

	if len(kafkaCfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers is not specified")
	}

	return &Config{
		Server: serverCfg,
		DB: dbCfg,
		JWT: jwtCfg,
		Kafka: kafkaCfg,
	}, nil
}
