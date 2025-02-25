package handlers

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsersHandler memungkinkan admin melihat semua user
func GetAllUsersHandler(c *gin.Context) {
	var users []models.User

	// Ambil semua user dari database
	if err := config.DB.Select("id, name, email, role, created_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
