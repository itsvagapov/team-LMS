package kafkabro

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(
	broker string,
	topic string,
	partitions int,
	replicationFactor int,
) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("dial broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(
		controller.Host,
		strconv.Itoa(controller.Port),
	)

	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return fmt.Errorf("create topic: %w", err)
	}

	return nil
}
