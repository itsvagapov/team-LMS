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

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(service service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) RegisterRoutes(router *gin.Engine) {
	courses := router.Group("/courses")
	{
		courses.GET("", middleware.RequireRoles("student", "teacher", "admin"), h.GetCourses)
		courses.GET("/:id", middleware.RequireRoles("student", "teacher", "admin"), h.GetCourseByID)
		courses.POST("", middleware.RequireRoles("teacher", "admin"), h.CreateCourse)
		courses.PATCH("/:id", middleware.RequireRoles("teacher", "admin"), h.UpdateCourse)
		courses.DELETE("/:id", middleware.RequireRoles("teacher", "admin"), h.DeleteCourse)
	}
}

func (h *CourseHandler) CreateCourse(ctx *gin.Context) {
	var req dto.CreateCourseRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	teacherID := ctx.GetUint(middleware.UserIDContextKey)

	course, err := h.service.CreateCourse(
		ctx.Request.Context(),
		teacherID,
		req.Title,
		req.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCourseTitle):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrInvalidCourseDescription):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create course",
			})
		}

		return
	}

	ctx.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourses(ctx *gin.Context) {
	courses, err := h.service.GetCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get courses",
		})
		return
	}

	ctx.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) GetCourseByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourseNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get course",
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, course)
}

func (h *CourseHandler) UpdateCourse(ctx *gin.Context) {
	idParam := ctx.Param("id")

	courseID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var req dto.UpdateCourseRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	actorID := ctx.GetUint(middleware.UserIDContextKey)
	actorRole := ctx.GetString(middleware.UserRoleContextKey)

	course, err := h.service.UpdateCourse(
		ctx.Request.Context(),
		actorID,
		actorRole,
		uint(courseID),
		req.Title,
		req.Description,
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

		case errors.Is(err, service.ErrInvalidCourseTitle),
			errors.Is(err, service.ErrInvalidCourseDescription),
			errors.Is(err, service.ErrNoFieldsToUpdate):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update course",
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, course)
}

func (h *CourseHandler) DeleteCourse(ctx *gin.Context) {
	idParam := ctx.Param("id")

	courseID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || courseID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	actorID := ctx.GetUint(middleware.UserIDContextKey)
	actorRole := ctx.GetString(middleware.UserRoleContextKey)

	err = h.service.DeleteCourse(
		ctx.Request.Context(),
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
				"error": "failed to delete course",
			})
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}
