package routes

import (
	"go-login-api/api/handlers"
	"go-login-api/api/handlers/auth"
	"go-login-api/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")

	// Public routes
	api.POST("/register", auth.RegisterHandler)
	api.POST("/login", auth.LoginHandler)
	api.POST("/refresh-token", auth.RefreshTokenHandler)
	api.POST("/forgot-password", auth.ForgotPasswordHandler)
	api.POST("/reset-password", auth.ResetPasswordHandler)

	// Protected routes (hanya bisa diakses dengan token JWT)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// User biasa bisa mengakses profil dan logout
	protected.GET("/profile", handlers.ProfileHandler)
	protected.POST("/logout", auth.LogoutHandler)
	protected.PUT("/update-profile", handlers.UpdateProfileHandler)
	protected.DELETE("/delete-account", auth.DeleteAccountHandler)

	// Grup khusus admin
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, Admin!"})
	})
	admin.POST("/register", auth.RegisterAdminHandler)
	admin.GET("/get-users", handlers.GetAllUsersHandler)
	admin.DELETE("/delete-user/:id", auth.DeleteUserByAdminHandler)

	return router
}
