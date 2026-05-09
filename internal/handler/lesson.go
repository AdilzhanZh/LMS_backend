package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateLesson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson id",
		})
		return
	}

	var lesson models.UpdateLesson
	if err = c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	updatedID, err := h.services.Lesson.Update(c.Request.Context(), id, lesson)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrLessonNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson to update not found"})
		case errors.Is(err, models.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": updatedID,
	})
}

func (h *Handler) CreateLesson(c *gin.Context) {
	var lesson models.CreateLesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	id, err := h.services.Lesson.Create(c.Request.Context(), lesson)
	if err != nil {
		if errors.Is(err, models.ErrCourseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a lesson"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *Handler) DeleteLesson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson id",
		})
		return
	}

	if err = h.services.Lesson.DeleteByID(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson to delete " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetLessonByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson id",
		})
		return
	}

	lesson, err := h.services.Lesson.GetLessonByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrLessonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Lesson " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *Handler) GetLessons(c *gin.Context) {
	lessons, err := h.services.Lesson.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get lessons"})
		return
	}

	c.JSON(http.StatusOK, lessons)
}
