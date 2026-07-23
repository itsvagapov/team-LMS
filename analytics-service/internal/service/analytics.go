package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"activity-analytics-service/internal/model"
	"activity-analytics-service/internal/repository"
)

type AnalyticsPublisher interface {
	Publish(ctx context.Context, eventType string, payload any) error
}

type AnalyticsService struct {
	events      *repository.ActivityEventRepository
	courses     *repository.CourseStatsRepository
	users       *repository.UserActivityStatsRepository
	publisher   AnalyticsPublisher
	nowFunction func() time.Time
}

type IncomingEvent struct {
	EventType     string
	SourceService string
	Payload       map[string]any
	RawPayload    []byte
	CreatedAt     time.Time
}

func NewAnalyticsService(
	events *repository.ActivityEventRepository,
	courses *repository.CourseStatsRepository,
	users *repository.UserActivityStatsRepository,
	publisher AnalyticsPublisher,
) *AnalyticsService {
	return &AnalyticsService{
		events:      events,
		courses:     courses,
		users:       users,
		publisher:   publisher,
		nowFunction: time.Now,
	}
}

func (s *AnalyticsService) HandleEvent(ctx context.Context, incoming IncomingEvent) error {
	if incoming.EventType == "" {
		return fmt.Errorf("event type is empty")
	}
	if incoming.CreatedAt.IsZero() {
		incoming.CreatedAt = s.nowFunction().UTC()
	}

	event := &model.ActivityEvent{
		EventType:     incoming.EventType,
		SourceService: incoming.SourceService,
		ActorID:       extractActorID(incoming.Payload),
		EntityType:    entityType(incoming.EventType),
		EntityID:      extractEntityID(incoming.EventType, incoming.Payload),
		Payload:       string(incoming.RawPayload),
		CreatedAt:     incoming.CreatedAt,
	}

	if err := s.events.Create(event); err != nil {
		return fmt.Errorf("save activity event: %w", err)
	}

	if err := s.updateStats(incoming); err != nil {
		return fmt.Errorf("update stats: %w", err)
	}

	if err := s.publisher.Publish(ctx, "activity.saved", map[string]any{
		"event_id":   event.ID,
		"event":      incoming.EventType,
		"created_at": s.nowFunction().UTC(),
	}); err != nil {
		log.Printf("publish activity.saved: %v", err)
	}

	return nil
}

func (s *AnalyticsService) ListEvents(limit, offset int) ([]model.ActivityEvent, error) {
	return s.events.List(normalizeLimit(limit), normalizeOffset(offset))
}

func (s *AnalyticsService) GetUserActivity(userID uint, limit, offset int) (*model.UserActivityStats, []model.ActivityEvent, error) {
	stats, err := s.users.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}
	events, err := s.events.ListByUser(userID, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return nil, nil, err
	}
	return stats, events, nil
}

func (s *AnalyticsService) GetCourseStats(courseID uint) (*model.CourseStats, error) {
	return s.courses.GetByCourseID(courseID)
}

func (s *AnalyticsService) GetDashboard() (*model.Dashboard, error) {
	eventsCount, err := s.events.Count()
	if err != nil {
		return nil, err
	}
	usersCount, err := s.users.Count()
	if err != nil {
		return nil, err
	}
	coursesCount, err := s.courses.Count()
	if err != nil {
		return nil, err
	}
	students, submissions, reviewed, err := s.courses.DashboardSums()
	if err != nil {
		return nil, err
	}

	return &model.Dashboard{
		EventsCount:              eventsCount,
		UsersCount:               usersCount,
		CoursesCount:             coursesCount,
		StudentsEnrolledCount:    students,
		SubmissionsCount:         submissions,
		ReviewedSubmissionsCount: reviewed,
	}, nil
}

func (s *AnalyticsService) updateStats(incoming IncomingEvent) error {
	courseID := numberFromPayload(incoming.Payload, "course_id")
	userID := numberFromPayload(incoming.Payload, "user_id")
	studentID := numberFromPayload(incoming.Payload, "student_id")

	switch incoming.EventType {
	case "user.registered":
		if userID != nil {
			if err := s.users.MarkRegistered(*userID, incoming.CreatedAt); err != nil {
				return err
			}
			s.publishUserStatsUpdated(*userID, "registered_at", incoming.CreatedAt)
			return nil
		}
	case "user.logged_in", "user.role_changed":
		if userID != nil {
			return s.touchUser(*userID, incoming.CreatedAt)
		}
	case "course.created":
		if courseID != nil {
			return s.courses.Ensure(*courseID)
		}
	case "lesson.created":
		if courseID != nil {
			return s.incrementCourse(*courseID, "lessons_count", incoming.CreatedAt)
		}
	case "student.enrolled":
		if courseID != nil {
			if err := s.incrementCourse(*courseID, "students_count", incoming.CreatedAt); err != nil {
				return err
			}
		}
		if studentID != nil {
			return s.incrementUser(*studentID, "courses_enrolled_count", incoming.CreatedAt)
		}
	case "homework.created":
		if courseID != nil {
			return s.incrementCourse(*courseID, "homeworks_count", incoming.CreatedAt)
		}
	case "submission.created":
		if courseID != nil {
			if err := s.incrementCourse(*courseID, "submissions_count", incoming.CreatedAt); err != nil {
				return err
			}
		}
		if studentID != nil {
			return s.incrementUser(*studentID, "submissions_count", incoming.CreatedAt)
		}
	case "submission.reviewed":
		if courseID != nil {
			return s.incrementCourse(*courseID, "reviewed_submissions_count", incoming.CreatedAt)
		}
		if studentID != nil {
			return s.touchUser(*studentID, incoming.CreatedAt)
		}
	}

	if actorID := extractActorID(incoming.Payload); actorID != nil {
		return s.touchUser(*actorID, incoming.CreatedAt)
	}
	return nil
}

func (s *AnalyticsService) incrementCourse(courseID uint, column string, createdAt time.Time) error {
	if err := s.courses.Increment(courseID, column, 1); err != nil {
		return err
	}
	if err := s.publisher.Publish(context.Background(), "course.stats.updated", map[string]any{
		"course_id":  courseID,
		"field":      column,
		"created_at": createdAt,
	}); err != nil {
		log.Printf("publish course.stats.updated: %v", err)
	}
	return nil
}

func (s *AnalyticsService) incrementUser(userID uint, column string, createdAt time.Time) error {
	if err := s.users.Increment(userID, column, 1, createdAt); err != nil {
		return err
	}
	s.publishUserStatsUpdated(userID, column, createdAt)
	return nil
}

func (s *AnalyticsService) touchUser(userID uint, createdAt time.Time) error {
	if err := s.users.Touch(userID, createdAt); err != nil {
		return err
	}
	s.publishUserStatsUpdated(userID, "last_activity_at", createdAt)
	return nil
}

func (s *AnalyticsService) publishUserStatsUpdated(userID uint, field string, createdAt time.Time) {
	if err := s.publisher.Publish(context.Background(), "user.stats.updated", map[string]any{
		"user_id":    userID,
		"field":      field,
		"created_at": createdAt,
	}); err != nil {
		log.Printf("publish user.stats.updated: %v", err)
	}
}

func ParseIncomingEvent(sourceService string, raw []byte) (IncomingEvent, error) {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return IncomingEvent{}, err
	}

	eventType, _ := payload["event"].(string)
	createdAt := parseTime(payload["created_at"])

	return IncomingEvent{
		EventType:     eventType,
		SourceService: sourceService,
		Payload:       payload,
		RawPayload:    raw,
		CreatedAt:     createdAt,
	}, nil
}

func parseTime(value any) time.Time {
	raw, ok := value.(string)
	if !ok || raw == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func extractActorID(payload map[string]any) *uint {
	for _, key := range []string{"user_id", "student_id", "teacher_id", "created_by", "reviewed_by", "changed_by"} {
		if value := numberFromPayload(payload, key); value != nil {
			return value
		}
	}
	return nil
}

func extractEntityID(eventType string, payload map[string]any) *uint {
	keysByPrefix := map[string][]string{
		"user.":       {"user_id"},
		"course.":     {"course_id"},
		"lesson.":     {"lesson_id"},
		"student.":    {"student_id"},
		"homework.":   {"homework_id"},
		"submission.": {"submission_id"},
	}
	for prefix, keys := range keysByPrefix {
		if len(eventType) >= len(prefix) && eventType[:len(prefix)] == prefix {
			for _, key := range keys {
				if value := numberFromPayload(payload, key); value != nil {
					return value
				}
			}
		}
	}
	return nil
}

func entityType(eventType string) string {
	for i, ch := range eventType {
		if ch == '.' {
			return eventType[:i]
		}
	}
	return "unknown"
}

func numberFromPayload(payload map[string]any, key string) *uint {
	value, ok := payload[key]
	if !ok {
		return nil
	}

	var parsed uint64
	switch typed := value.(type) {
	case float64:
		if typed <= 0 {
			return nil
		}
		parsed = uint64(typed)
	case string:
		value, err := strconv.ParseUint(typed, 10, 64)
		if err != nil || value == 0 {
			return nil
		}
		parsed = value
	default:
		return nil
	}

	result := uint(parsed)
	return &result
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 200 {
		return 200
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
