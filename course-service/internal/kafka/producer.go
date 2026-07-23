package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

const CourseEventsTopic = "courses.events"

type Producer interface {
	Publish(ctx context.Context, event any) error
	Close() error
}

type producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string) Producer {
	return &producer{
		writer: &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Topic:                  CourseEventsTopic,
			Balancer:               &kafkago.LeastBytes{},
			RequiredAcks:           kafkago.RequireAll,
			WriteTimeout:           10 * time.Second,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *producer) Publish(
	ctx context.Context,
	event any,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal kafka event: %w", err)
	}

	if err := p.writer.WriteMessages(
		ctx,
		kafkago.Message{
			Value: payload,
			Time:  time.Now().UTC(),
		},
	); err != nil {
		return fmt.Errorf("publish kafka event: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
