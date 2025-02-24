package middleware

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole memastikan hanya user dengan role tertentu yang bisa mengakses endpoint
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil user ID dari context (di-set oleh middleware JWT)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Ambil data user dari database
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Cek apakah user memiliki peran yang sesuai
		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
