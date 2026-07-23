package kafka

import (
	"context"
	"errors"
	"log"
	"strings"

	"activity-analytics-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	readers []*kafka.Reader
	service *service.AnalyticsService
}

func NewConsumer(brokers []string, groupID string, topics []string, service *service.AnalyticsService) *Consumer {
	readers := make([]*kafka.Reader, 0, len(topics))
	for _, topic := range topics {
		readers = append(readers, kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}))
	}
	return &Consumer{readers: readers, service: service}
}

func (c *Consumer) Start(ctx context.Context) {
	for _, reader := range c.readers {
		go c.consumeTopic(ctx, reader)
	}
}

func (c *Consumer) consumeTopic(ctx context.Context, reader *kafka.Reader) {
	topic := reader.Config().Topic
	sourceService := sourceFromTopic(topic)
	log.Printf("kafka consumer started for topic %s", topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("kafka consumer error for topic %s: %v", topic, err)
			continue
		}

		event, err := service.ParseIncomingEvent(sourceService, msg.Value)
		if err != nil {
			log.Printf("parse kafka event from topic %s: %v", topic, err)
			continue
		}

		if err := c.service.HandleEvent(ctx, event); err != nil {
			log.Printf("handle kafka event %s from topic %s: %v", event.EventType, topic, err)
			continue
		}

		log.Printf("processed kafka event %s from topic %s", event.EventType, topic)
	}
}

func (c *Consumer) Close() error {
	var closeErr error
	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			closeErr = err
		}
	}
	return closeErr
}

func sourceFromTopic(topic string) string {
	switch {
	case strings.HasPrefix(topic, "users."):
		return "user-service"
	case strings.HasPrefix(topic, "courses."):
		return "course-service"
	case strings.HasPrefix(topic, "homework."):
		return "homework-service"
	default:
		return topic
	}
}
