package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/itsvagapov/team-LMS/course-service/internal/kafka"
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"github.com/itsvagapov/team-LMS/course-service/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidLessonTitle   = errors.New("lesson title must contain from 1 to 50 characters")
	ErrInvalidLessonContent = errors.New("lesson content must contain from 1 to 2000 characters")
	ErrLessonNotFound       = errors.New("lesson not found")
)

type LessonService interface {
	CreateLesson(ctx context.Context, actorID uint, actorRole string, courseID uint, title string, content string) (*models.Lesson, error)
	GetLessons(actorID uint, actorRole string, courseID uint) ([]models.Lesson, error)
	GetLessonByID(id uint) (*models.Lesson, error)
	PublishLesson(ctx context.Context, actorID uint, actorRole string, lessonID uint) (*models.Lesson, error)
}

type lessonService struct {
	courseService CourseService
	lessonRepo    repository.LessonRepository
	producer      kafka.Producer
}

func NewLessonService(courseService CourseService, lessonRepo repository.LessonRepository, producer kafka.Producer) LessonService {
	return &lessonService{
		courseService: courseService,
		lessonRepo:    lessonRepo,
		producer:      producer,
	}
}

func (s *lessonService) CreateLesson(ctx context.Context, actorID uint, actorRole string, courseID uint, title string, content string) (*models.Lesson, error) {
	course, err := s.courseService.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	if actorRole != "admin" {
		if actorRole != "teacher" || course.TeacherID != actorID {
			return nil, ErrAccessDenied
		}
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || utf8.RuneCountInString(title) > 50 {
		return nil, ErrInvalidLessonTitle
	}

	if content == "" || utf8.RuneCountInString(content) > 2000 {
		return nil, ErrInvalidLessonContent
	}

	orderNumber, err := s.lessonRepo.GetNextOrderNumber(courseID)
	if err != nil {
		return nil, err
	}

	lesson := &models.Lesson{
		CourseID:    courseID,
		Title:       title,
		Content:     content,
		OrderNumber: orderNumber,
		IsPublished: false,
	}

	if err := s.lessonRepo.Create(lesson); err != nil {
		return nil, err
	}

	event := kafka.LessonCreatedEvent{
		Event:       kafka.EventLessonCreated,
		LessonID:    lesson.ID,
		CourseID:    lesson.CourseID,
		Title:       lesson.Title,
		OrderNumber: lesson.OrderNumber,
		CreatedAt:   lesson.CreatedAt.UTC(),
	}

	if err := s.producer.Publish(ctx, event); err != nil {
		log.Printf(
			"failed to publish %s event: %v",
			kafka.EventLessonCreated,
			err,
		)
	}
	log.Printf(
		"lesson created: lesson_id=%d course_id=%d actor_id=%d actor_role=%s order_number=%d",
		lesson.ID,
		lesson.CourseID,
		actorID,
		actorRole,
		lesson.OrderNumber,
	)

	return lesson, nil
}

func (s *lessonService) GetLessons(actorID uint, actorRole string, courseID uint) ([]models.Lesson, error) {
	course, err := s.courseService.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}
	publishedOnly := true

	switch actorRole {
	case "admin":
		publishedOnly = false

	case "teacher":
		if course.TeacherID == actorID {
			publishedOnly = false
		}

	case "student":
		publishedOnly = true

	default:
		return nil, ErrAccessDenied
	}

	return s.lessonRepo.FindByCourseID(courseID, publishedOnly)
}

func (s *lessonService) GetLessonByID(id uint) (*models.Lesson, error) {
	lesson, err := s.lessonRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLessonNotFound
		}

		return nil, err
	}

	return lesson, nil
}

func (s *lessonService) PublishLesson(ctx context.Context, actorID uint, actorRole string, lessonID uint) (*models.Lesson, error) {
	lesson, err := s.GetLessonByID(lessonID)
	if err != nil {
		return nil, err
	}

	course, err := s.courseService.GetCourseByID(lesson.CourseID)
	if err != nil {
		return nil, err
	}

	if actorRole != "admin" {
		if actorRole != "teacher" || course.TeacherID != actorID {
			return nil, ErrAccessDenied
		}
	}

	if lesson.IsPublished {
		return lesson, nil
	}

	lesson.IsPublished = true

	if err := s.lessonRepo.Update(lesson); err != nil {
		return nil, err
	}

	event := kafka.LessonPublishedEvent{
		Event:       kafka.EventLessonPublished,
		LessonID:    lesson.ID,
		CourseID:    lesson.CourseID,
		IsPublished: lesson.IsPublished,
		CreatedAt:   lesson.UpdatedAt.UTC(),
	}

	if err := s.producer.Publish(ctx, event); err != nil {
		log.Printf(
			"failed to publish %s event: %v",
			kafka.EventLessonPublished,
			err,
		)
	}
	log.Printf(
		"lesson published: lesson_id=%d course_id=%d actor_id=%d actor_role=%s",
		lesson.ID,
		lesson.CourseID,
		actorID,
		actorRole,
	)

	return lesson, nil
}
