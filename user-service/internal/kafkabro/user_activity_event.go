package kafkabro

import (
	"time"

	"github.com/itsvagapov/team-LMS/user-service/internal/model"
)

type UserActivityEventMessage struct {
	UserID uint `json:"user_id"`

	EventType UserActivityEventType `json:"event_type"`

	SourceService SourceService `json:"source_service"`

	Payload string `json:"payload"`

	CreatedAt time.Time `json:"created_at"`
}

type UserLoggedInPayload struct {
	Role model.UserRole `json:"role"`
}
