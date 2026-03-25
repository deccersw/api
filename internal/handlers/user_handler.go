package handlers

import (
	"net/http"
	"todo_api/internal/domain"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var RegisterRequest RegisterRequest
		if err := c.BindJSON(&RegisterRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(RegisterRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RegisterRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash passwrord" + err.Error()})
			return
		}

		user := &domain.User{
			Email:    RegisterRequest.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateUser(pool, user)
		if err != nil {
			if err.Error() != "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email alredy register" + err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdUser)
	}
}
