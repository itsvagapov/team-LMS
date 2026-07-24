package kafkabro

import (
	"fmt"
	"net"
	"strconv"

	"github.com/itsvagapov/team-LMS/user-service/internal/config"
	"github.com/segmentio/kafka-go"
)

const (
	UsersEventsTopic string = "users.events"
)

func CreateTopic(cfg config.KafkaConfig, topic string) error {
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		return fmt.Errorf("dial broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerConn, err := kafka.Dial(
		"tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	)
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	if err := controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     cfg.PartitionsNum,
		ReplicationFactor: cfg.ReplicationFactor,
	}); err != nil {
		return fmt.Errorf("create topic: %w", err)
	}

	return nil
}
