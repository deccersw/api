// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"
	"todo_api/internal/repository"
	"todo_api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := database.Connect(cfg.Port, cfg.Host, cfg.DBName, cfg.SSlmode, cfg.User)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	todoRepo := repository.NewTodoRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	todoService := service.NewTodoService(todoRepo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret, 24*time.Hour)

	todoHandler := handlers.NewTodoHandler(todoService)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	todo := router.Group("/todo")
	todo.Use(middleware.AuthMiddleware(cfg))
	{
		todo.POST("", todoHandler.Create)
		todo.GET("", todoHandler.GetAll)
		todo.GET("/:id", todoHandler.GetByID)
		todo.PATCH("/:id", todoHandler.Update)
		todo.DELETE("/:id", todoHandler.Delete)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: router,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")

}
