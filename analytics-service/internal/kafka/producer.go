package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(eventType),
		Value: body,
		Time:  time.Now().UTC(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("kafka producer error: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
