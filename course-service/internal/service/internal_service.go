package service

import (
	"github.com/itsvagapov/team-LMS/course-service/internal/models"
	"github.com/itsvagapov/team-LMS/course-service/internal/repository"
)

type InternalService interface {
	GetCourse(courseID uint) (*models.Course, error)
	GetLesson(lessonID uint) (*models.Lesson, error)
	IsStudentEnrolled(courseID uint, studentID uint) (bool, error)
}

type internalService struct {
	courseService     CourseService
	lessonService     LessonService
	courseStudentRepo repository.CourseStudentRepository
}

func NewInternalService(courseService CourseService, lessonService LessonService, courseStudentRepo repository.CourseStudentRepository) InternalService {
	return &internalService{
		courseService:     courseService,
		lessonService:     lessonService,
		courseStudentRepo: courseStudentRepo,
	}
}

func (s *internalService) GetCourse(courseID uint) (*models.Course, error) {
	return s.courseService.GetCourseByID(courseID)
}

func (s *internalService) GetLesson(lessonID uint) (*models.Lesson, error) {
	return s.lessonService.GetLessonByID(lessonID)
}

func (s *internalService) IsStudentEnrolled(courseID uint, studentID uint) (bool, error) {
	if _, err := s.courseService.GetCourseByID(courseID); err != nil {
		return false, err
	}

	return s.courseStudentRepo.Exists(courseID, studentID)
}
