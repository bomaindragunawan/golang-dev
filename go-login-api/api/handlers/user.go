package handlers

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler untuk mendapatkan profil user yang sedang login
func ProfileHandler(c *gin.Context) {
	// Ambil email dari context (yang diset di middleware)
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return data user
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
