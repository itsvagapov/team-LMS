package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/course-service/internal/middleware"
	"github.com/itsvagapov/team-LMS/course-service/internal/service"
)

type EnrollmentHandler struct {
	service service.EnrollmentService
}

func NewEnrollmentHandler(service service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

func (h *EnrollmentHandler) RegisterRoutes(router *gin.Engine) {
	courses := router.Group("/courses")
	{
		courses.POST("/:id/enroll", middleware.RequireRoles("student"), h.EnrollStudent)
		courses.GET("/:id/students", middleware.RequireRoles("teacher", "admin"), h.GetCourseStudents)
	}
}

func (h *EnrollmentHandler) EnrollStudent(ctx *gin.Context) {
	courseIDParam := ctx.Param("id")

	courseID, err := strconv.ParseUint(courseIDParam, 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	studentID := ctx.GetUint(middleware.UserIDContextKey)
	studentRole := ctx.GetString(middleware.UserRoleContextKey)

	courseStudent, err := h.service.EnrollStudent(
		ctx.Request.Context(),
		studentID,
		studentRole,
		uint(courseID),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourseNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrAccessDenied):
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrAlreadyEnrolled):
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to enroll student",
			})
		}

		return
	}

	ctx.JSON(http.StatusCreated, courseStudent)
}

func (h *EnrollmentHandler) GetCourseStudents(ctx *gin.Context) {
	courseIDParam := ctx.Param("id")

	courseID, err := strconv.ParseUint(courseIDParam, 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	actorID := ctx.GetUint(middleware.UserIDContextKey)
	actorRole := ctx.GetString(middleware.UserRoleContextKey)

	courseStudents, err := h.service.GetCourseStudents(
		actorID,
		actorRole,
		uint(courseID),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourseNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrAccessDenied):
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get course students",
			})
		}

		return
	}

	ctx.JSON(http.StatusOK, courseStudents)
}
