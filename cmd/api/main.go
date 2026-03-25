package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.Port, cfg.Host, cfg.DBName, cfg.SSlmode, cfg.User)

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":            "TODO API is running",
			"status":             "success",
			"database_connected": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))
	protected := router.Group("/todo")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", handlers.CreateTodoHandler(pool))
		protected.GET("", handlers.GetAllTodoHandler(pool))
		protected.GET("/:id", handlers.GetTodoByIdHandler(pool))
		protected.PATCH("/:id", handlers.UpdateTodoHandler(pool))
		protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))
	}

	router.Run(":" + cfg.PORT)

}
