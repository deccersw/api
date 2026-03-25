package handlers

import (
	"net/http"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/domain"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
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

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var LoginRequest LoginRequest

		if err := c.BindJSON(&LoginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, LoginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials" + err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials " + err.Error()})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed generate token" + err.Error()})
		}
		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}
