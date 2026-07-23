package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/course-service/internal/service"
)

type InternalHandler struct {
	service service.InternalService
}

func NewInternalHandler(service service.InternalService) *InternalHandler {
	return &InternalHandler{service: service}
}

func (h *InternalHandler) RegisterRoutes(router *gin.Engine) {
	internal := router.Group("/internal")
	{
		internal.GET("/courses/:id", h.GetCourse)
		internal.GET("/lessons/:id", h.GetLesson)
		internal.GET("/courses/:id/students/:student_id", h.IsStudentEnrolled)
	}
}

func (h *InternalHandler) GetCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	course, err := h.service.GetCourse(uint(courseID))
	if err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get course",
		})
		return
	}

	ctx.JSON(http.StatusOK, course)
}

func (h *InternalHandler) GetLesson(ctx *gin.Context) {
	lessonID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || lessonID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson id",
		})
		return
	}

	lesson, err := h.service.GetLesson(uint(lessonID))
	if err != nil {
		if errors.Is(err, service.ErrLessonNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get lesson",
		})
		return
	}

	ctx.JSON(http.StatusOK, lesson)
}

func (h *InternalHandler) IsStudentEnrolled(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(
		ctx.Param("id"), 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	studentID, err := strconv.ParseUint(
		ctx.Param("student_id"), 10, 64)
	if err != nil || studentID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student id",
		})
		return
	}

	enrolled, err := h.service.IsStudentEnrolled(
		uint(courseID),
		uint(studentID),
	)
	if err != nil {
		if errors.Is(err, service.ErrCourseNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check student enrollment",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"course_id":  courseID,
		"student_id": studentID,
		"enrolled":   enrolled,
	})
}
