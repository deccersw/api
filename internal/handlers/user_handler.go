package handlers

import (
	"context"
	"errors"
	"net/http"
	"todo_api/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Register(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), domain.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
			return
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token})
}
