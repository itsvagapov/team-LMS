package service

import (
	"context"
	"errors"
	"log"

	"github.com/itsvagapov/team-LMS/course-service/internal/kafka"
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"github.com/itsvagapov/team-LMS/course-service/internal/repository"
)

var ErrAlreadyEnrolled = errors.New(
	"student is already enrolled in this course",
)

type EnrollmentService interface {
	EnrollStudent(ctx context.Context, studentID uint, studentRole string, courseID uint) (*models.CourseStudent, error)
	GetCourseStudents(actorID uint, actorRole string, courseID uint) ([]models.CourseStudent, error)
}

type enrollmentService struct {
	courseService     CourseService
	courseStudentRepo repository.CourseStudentRepository
	producer          kafka.Producer
}

func NewEnrollmentService(courseService CourseService, courseStudentRepo repository.CourseStudentRepository, producer kafka.Producer) EnrollmentService {
	return &enrollmentService{
		courseService:     courseService,
		courseStudentRepo: courseStudentRepo,
		producer:          producer,
	}
}

func (s *enrollmentService) EnrollStudent(ctx context.Context, studentID uint, studentRole string, courseID uint) (*models.CourseStudent, error) {
	if studentRole != "student" {
		return nil, ErrAccessDenied
	}

	if _, err := s.courseService.GetCourseByID(courseID); err != nil {
		return nil, err
	}

	exists, err := s.courseStudentRepo.Exists(
		courseID,
		studentID,
	)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrAlreadyEnrolled
	}

	courseStudent := &models.CourseStudent{
		CourseID:  courseID,
		StudentID: studentID,
	}

	if err := s.courseStudentRepo.Create(courseStudent); err != nil {
		return nil, err
	}

	event := kafka.StudentEnrolledEvent{
		Event:     kafka.EventStudentEnrolled,
		CourseID:  courseStudent.CourseID,
		StudentID: courseStudent.StudentID,
		CreatedAt: courseStudent.CreatedAt.UTC(),
	}

	if err := s.producer.Publish(ctx, event); err != nil {
		log.Printf(
			"failed to publish %s event: %v",
			kafka.EventStudentEnrolled,
			err,
		)
	}

	log.Printf(
		"student enrolled: student_id=%d course_id=%d",
		courseStudent.StudentID,
		courseStudent.CourseID,
	)

	return courseStudent, nil
}

func (s *enrollmentService) GetCourseStudents(actorID uint, actorRole string, courseID uint) ([]models.CourseStudent, error) {
	course, err := s.courseService.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	if actorRole != "admin" {
		if actorRole != "teacher" || course.TeacherID != actorID {
			return nil, ErrAccessDenied
		}
	}

	return s.courseStudentRepo.FindByCourseID(courseID)
}
