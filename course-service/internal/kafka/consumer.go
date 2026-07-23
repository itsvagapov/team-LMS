package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"github.com/itsvagapov/team-LMS/course-service/internal/repository"
	kafkago "github.com/segmentio/kafka-go"
)

const HomeworkEventsTopic = "homework.events"

type Consumer struct {
	reader *kafkago.Reader
	repo   repository.CourseEventRepository
}

type eventEnvelope struct {
	Event string `json:"event"`
}

func NewConsumer(
	brokers []string,
	groupID string,
	repo repository.CourseEventRepository,
) *Consumer {
	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    HomeworkEventsTopic,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		repo: repo,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Printf("failed to fetch kafka message: %v", err)
			continue
		}

		var envelope eventEnvelope

		if err := json.Unmarshal(message.Value, &envelope); err != nil {
			log.Printf("failed to decode kafka event: %v", err)
			c.commitMessage(ctx, message)
			continue
		}

		if envelope.Event == "" {
			log.Printf("received kafka event without event type")
			c.commitMessage(ctx, message)
			continue
		}

		event := &models.CourseEvent{
			EventType: envelope.Event,
			Payload:   string(message.Value),
			CreatedAt: time.Now().UTC(),
		}

		if err := c.repo.Create(event); err != nil {
			log.Printf(
				"failed to save kafka event %s: %v",
				envelope.Event,
				err,
			)
			continue
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			log.Printf(
				"failed to commit kafka event %s: %v",
				envelope.Event,
				err,
			)
			continue
		}

		log.Printf(
			"kafka event consumed: topic=%s event=%s",
			message.Topic,
			envelope.Event,
		)
	}
}

func (c *Consumer) commitMessage(
	ctx context.Context,
	message kafkago.Message,
) {
	if err := c.reader.CommitMessages(ctx, message); err != nil {
		log.Printf("failed to commit invalid kafka message: %v", err)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
