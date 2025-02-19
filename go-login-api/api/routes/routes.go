package routes

import (
	"go-login-api/api/handlers"
	"go-login-api/api/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)

	// Protected routes (hanya bisa diakses dengan token JWT)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", handlers.ProfileHandler)

	return router
}
