package kafkabro

import (
	"context"

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
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) Send(ctx context.Context, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}