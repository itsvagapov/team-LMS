package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"unicode/utf8"

	coursekafka "github.com/itsvagapov/team-LMS/course-service/internal/kafka"
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"github.com/itsvagapov/team-LMS/course-service/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidCourseTitle       = errors.New("course title must contain from 1 to 50 characters")
	ErrInvalidCourseDescription = errors.New("course description must not exceed 200 characters")
	ErrCourseNotFound           = errors.New("course not found")
	ErrAccessDenied             = errors.New("access denied")
	ErrNoFieldsToUpdate         = errors.New("no fields to update")
)

type CourseService interface {
	CreateCourse(ctx context.Context, teacherID uint, title, description string) (*models.Course, error)
	GetCourses() ([]models.Course, error)
	GetCourseByID(id uint) (*models.Course, error)
	UpdateCourse(ctx context.Context, actorID uint, actorRole string, courseID uint, title *string, description *string) (*models.Course, error)
	DeleteCourse(ctx context.Context, actorID uint, actorRole string, courseID uint) error
}

type courseService struct {
	repo     repository.CourseRepository
	producer coursekafka.Producer
}

func NewCourseService(repo repository.CourseRepository, producer coursekafka.Producer) CourseService {
	return &courseService{
		repo:     repo,
		producer: producer,
	}
}

func (s *courseService) CreateCourse(ctx context.Context, teacherID uint, title string, description string) (*models.Course, error) {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" || utf8.RuneCountInString(title) > 50 {
		return nil, ErrInvalidCourseTitle
	}

	if utf8.RuneCountInString(description) > 200 {
		return nil, ErrInvalidCourseDescription
	}
	course := &models.Course{
		Title:       title,
		Description: description,
		TeacherID:   teacherID,
		IsActive:    true,
	}

	if err := s.repo.Create(course); err != nil {
		return nil, err
	}

	event := coursekafka.CourseCreatedEvent{
		Event:     coursekafka.EventCourseCreated,
		CourseID:  course.ID,
		TeacherID: course.TeacherID,
		Title:     course.Title,
		CreatedAt: course.CreatedAt.UTC(),
	}

	if err := s.producer.Publish(ctx, event); err != nil {
		log.Printf(
			"failed to publish %s event: %v",
			coursekafka.EventCourseCreated,
			err,
		)
	}
	log.Printf(
		"course created: course_id=%d teacher_id=%d",
		course.ID,
		course.TeacherID,
	)

	return course, nil
}

func (s *courseService) GetCourses() ([]models.Course, error) {
	return s.repo.FindAll()
}

func (s *courseService) GetCourseByID(id uint) (*models.Course, error) {
	course, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	return course, nil
}

func (s *courseService) UpdateCourse(ctx context.Context, actorID uint, actorRole string, courseID uint, title *string, description *string) (*models.Course, error) {
	course, err := s.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	if actorRole != "admin" {
		if actorRole != "teacher" || course.TeacherID != actorID {
			return nil, ErrAccessDenied
		}
	}

	if title == nil && description == nil {
		return nil, ErrNoFieldsToUpdate
	}

	if title != nil {
		normalizedTitle := strings.TrimSpace(*title)

		if normalizedTitle == "" || utf8.RuneCountInString(normalizedTitle) > 50 {
			return nil, ErrInvalidCourseTitle
		}
		course.Title = normalizedTitle
	}

	if description != nil {
		normalizedDescription := strings.TrimSpace(*description)

		if utf8.RuneCountInString(normalizedDescription) > 200 {
			return nil, ErrInvalidCourseDescription
		}
		course.Description = normalizedDescription
	}

	if err := s.repo.Update(course); err != nil {
		return nil, err
	}

	event := coursekafka.CourseUpdatedEvent{
		Event:       coursekafka.EventCourseUpdated,
		CourseID:    course.ID,
		TeacherID:   course.TeacherID,
		Title:       course.Title,
		Description: course.Description,
		IsActive:    course.IsActive,
		CreatedAt:   course.UpdatedAt.UTC(),
	}

	if err := s.producer.Publish(ctx, event); err != nil {
		log.Printf(
			"failed to publish %s event: %v",
			coursekafka.EventCourseUpdated,
			err,
		)
	}

	log.Printf(
		"course updated: course_id=%d actor_id=%d actor_role=%s",
		course.ID,
		actorID,
		actorRole,
	)
	return course, nil
}

func (s *courseService) DeleteCourse(ctx context.Context, actorID uint, actorRole string, courseID uint) error {
	course, err := s.GetCourseByID(courseID)
	if err != nil {
		return err
	}

	if actorRole != "admin" {
		if actorRole != "teacher" || course.TeacherID != actorID {
			return ErrAccessDenied
		}
	}

	course.IsActive = false

	if err := s.repo.Update(course); err != nil {
		return err
	}

	event := coursekafka.CourseDeletedEvent{
		Event:     coursekafka.EventCourseDeleted,
		CourseID:  course.ID,
		TeacherID: course.TeacherID,
		CreatedAt: course.UpdatedAt.UTC(),
	}

	if err := s.producer.Publish(ctx, event); err != nil {
		log.Printf(
			"failed to publish %s event: %v",
			coursekafka.EventCourseDeleted,
			err,
		)
	}

	log.Printf(
		"course deleted: course_id=%d actor_id=%d actor_role=%s",
		course.ID,
		actorID,
		actorRole,
	)

	return nil
}
