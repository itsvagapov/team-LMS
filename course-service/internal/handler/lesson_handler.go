package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/team-LMS/course-service/internal/dto"
	"github.com/itsvagapov/team-LMS/course-service/internal/middleware"
	"github.com/itsvagapov/team-LMS/course-service/internal/service"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(service service.LessonService) *LessonHandler {
	return &LessonHandler{service: service}
}

func (h *LessonHandler) RegisterRoutes(router *gin.Engine) {
	courses := router.Group("/courses")
	{
		courses.GET("/:id/lessons", middleware.RequireRoles("student", "teacher", "admin"), h.GetLessons)
		courses.POST("/:id/lessons", middleware.RequireRoles("teacher", "admin"), h.CreateLesson)
	}
	lessons := router.Group("/lessons")
	{
		lessons.PATCH("/:id/publish", middleware.RequireRoles("teacher", "admin"), h.PublishLesson)
	}
}

func (h *LessonHandler) GetLessons(ctx *gin.Context) {
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

	lessons, err := h.service.GetLessons(
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
				"error": "failed to get lessons",
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, lessons)
}

func (h *LessonHandler) CreateLesson(ctx *gin.Context) {
	courseIDParam := ctx.Param("id")

	courseID, err := strconv.ParseUint(courseIDParam, 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var req dto.CreateLessonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	actorID := ctx.GetUint(middleware.UserIDContextKey)
	actorRole := ctx.GetString(middleware.UserRoleContextKey)

	lesson, err := h.service.CreateLesson(
		ctx.Request.Context(),
		actorID,
		actorRole,
		uint(courseID),
		req.Title,
		req.Content,
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

		case errors.Is(err, service.ErrInvalidLessonTitle),
			errors.Is(err, service.ErrInvalidLessonContent):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create lesson",
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, lesson)
}

func (h *LessonHandler) PublishLesson(ctx *gin.Context) {
	lessonIDParam := ctx.Param("id")

	lessonID, err := strconv.ParseUint(lessonIDParam, 10, 64)
	if err != nil || lessonID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson id",
		})
		return
	}

	actorID := ctx.GetUint(middleware.UserIDContextKey)
	actorRole := ctx.GetString(middleware.UserRoleContextKey)

	lesson, err := h.service.PublishLesson(
		ctx.Request.Context(),
		actorID,
		actorRole,
		uint(lessonID),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrLessonNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
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
				"error": "failed to publish lesson",
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, lesson)

}
