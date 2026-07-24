package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	Brokers           []string
	PartitionsNum     int
	ReplicationFactor int
}

func New() (*Config, error) {
	is_docker := os.Getenv("IS_DOCKER")
	if is_docker == "" {
		err := godotenv.Load(".env.local")
		if err != nil {
			return nil, fmt.Errorf("failed to load env file: %w", err)
		}
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
		TTL:    ttl,
	}

	if jwtCfg.Secret == "" {
		return nil, errors.New("jwt secret is not specified")
	}

	partitionsNum, err := strconv.Atoi(os.Getenv("KAFKA_PARTITIONS_NUM"))
	if err != nil {
		return nil, errors.New("kafka partitions num is not specified")
	}

	replicationFactor, err := strconv.Atoi(os.Getenv("KAFKA_REPLICATION_FACTOR"))
	if err != nil {
		return nil, errors.New("kafka replication factor is not specified")
	}

	kafkaCfg := KafkaConfig{
		Brokers:           []string{os.Getenv("KAFKA_BROKERS")},
		PartitionsNum:     partitionsNum,
		ReplicationFactor: replicationFactor,
	}

	if len(kafkaCfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers is not specified")
	}

	return &Config{
		Server: serverCfg,
		DB:     dbCfg,
		JWT:    jwtCfg,
		Kafka:  kafkaCfg,
	}, nil
}

func NewDBConn(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
