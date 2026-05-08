package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	courseService *service.CourseService
	//TODO add other service
}

func NewHandler(cs *service.CourseService) *Handler {
	return &Handler{
		courseService: cs,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()

	r.GET("/courses", h.GetCourses)
	r.GET("/courses/:id", h.GetCourseByID)
	r.DELETE("/courses/:id", h.DeleteCourse)
	r.POST("/courses", h.CreateCourse)

	return r, nil
}

func (h *Handler) CreateCourse(c *gin.Context) {
	var course models.CreateCourse
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	id, err := h.courseService.Create(c.Request.Context(), course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to created course",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *Handler) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	if err = h.courseService.DeleteByID(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course to delete " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetCourseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	course, err := h.courseService.GetCourseByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Course " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.courseService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}
