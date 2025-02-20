package routes

import (
	"go-login-api/api/handlers"
	"go-login-api/api/handlers/auth"
	"go-login-api/api/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	api.POST("/register", auth.RegisterHandler)
	api.POST("/login", auth.LoginHandler)
	api.POST("/refresh-token", auth.RefreshTokenHandler)
	api.POST("/forgot-password", auth.ForgotPasswordHandler)
	api.POST("/reset-password", auth.ResetPasswordHandler)

	// Protected routes (hanya bisa diakses dengan token JWT)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", handlers.ProfileHandler)
	protected.POST("/logout", auth.LogoutHandler)

	return router
}
