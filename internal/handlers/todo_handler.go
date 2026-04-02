package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"todo_api/internal/domain"

	"github.com/gin-gonic/gin"
)

type TodoService interface {
	Create(ctx context.Context, input domain.CreateTodoInput, userID string) (*domain.Todo, error)
	GetAll(ctx context.Context, userID string) ([]domain.Todo, error)
	GetByID(ctx context.Context, id int, userID string) (*domain.Todo, error)
	Update(ctx context.Context, id int, input domain.UpdateTodoInput, userID string) (*domain.Todo, error)
	Delete(ctx context.Context, id int, userID string) error
}

type TodoHandler struct {
	service TodoService
}

func NewTodoHandler(service TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

type createTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type updateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
		return "", false
	}
	return userID.(string), true
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.Create(c.Request.Context(), domain.CreateTodoInput{
		Title:     req.Title,
		Completed: req.Completed,
	}, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	todos, err := h.service.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	todo, err := h.service.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title == nil && req.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field must be provided"})
		return
	}

	todo, err := h.service.Update(c.Request.Context(), id, domain.UpdateTodoInput{
		Title:     req.Title,
		Completed: req.Completed,
	}, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
}
