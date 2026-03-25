package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"

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
	router.POST("/todo", handlers.CreateTodoHandler(pool))
	router.GET("/todo", handlers.GetAllTodoHandler(pool))
	router.GET("/todo/:id", handlers.GetTodoByIdHandler(pool))
	router.PATCH("/todo/:id", handlers.UpdateTodoHandler(pool))
	router.DELETE("/todo/:id", handlers.DeleteTodoHandler(pool))
	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.Run(":" + cfg.PORT)

}
