package handler

import (
	"errors"
	"net/http"
	"strconv"

	"activity-analytics-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(service *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) ListEvents(c *gin.Context) {
	events, err := h.service.ListEvents(queryInt(c, "limit", 50), queryInt(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load activity events"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": events})
}

func (h *AnalyticsHandler) GetUserActivity(c *gin.Context) {
	userID, ok := pathUint(c, "id")
	if !ok {
		return
	}

	stats, events, err := h.service.GetUserActivity(userID, queryInt(c, "limit", 50), queryInt(c, "offset", 0))
	if err != nil {
		writeRepositoryError(c, err, "user activity not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":  stats,
		"events": events,
	})
}

func (h *AnalyticsHandler) GetCourseStats(c *gin.Context) {
	courseID, ok := pathUint(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetCourseStats(courseID)
	if err != nil {
		writeRepositoryError(c, err, "course stats not found")
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	dashboard, err := h.service.GetDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard"})
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

func pathUint(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(value), true
}

func queryInt(c *gin.Context, name string, fallback int) int {
	raw := c.Query(name)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func writeRepositoryError(c *gin.Context, err error, notFoundMessage string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundMessage})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
