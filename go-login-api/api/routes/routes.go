package routes

import (
	"go-login-api/api/handlers"
	"go-login-api/api/handlers/auth"
	"go-login-api/api/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()
	router.POST("/register", auth.RegisterHandler)
	router.POST("/login", auth.LoginHandler)
	router.POST("/refresh-token", auth.RefreshTokenHandler)

	// Protected routes (hanya bisa diakses dengan token JWT)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", handlers.ProfileHandler)
	protected.POST("/logout", auth.LogoutHandler)

	return router
}
