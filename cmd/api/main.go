package main

import (
	"log"
	"net/http"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// loading application configuration for environment variable
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to PostgresSQL database using connection string
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// closing database connection when application exits
	defer pool.Close()

	// creating gin router
	router := gin.Default()

	// Do not trust proxy headers (e.g. X-Forwarded-For); use the direct client connection IP.
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		// map[string]interface{}
		// map[string]any{}
		c.JSON(http.StatusOK, gin.H{
			"message":  "Todo API is running!",
			"status":   "Success",
			"database": "Connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", handlers.CreateTodoHandler(pool))
		protected.GET("", handlers.GetAllTodosHandler(pool))
		protected.GET("/:id", handlers.GetToDoByIDHandler(pool))
		protected.PUT("/:id", handlers.UpdateToDoHandler(pool))
		protected.DELETE("/:id", handlers.DeleteToDoHandler(pool))
	}

	// Middleware Test Route
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())

	// start http server on port PORT
	router.Run(":" + cfg.Port)
}
