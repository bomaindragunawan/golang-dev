package auth

import (
	"go-login-api/config"
	"go-login-api/helper"
	"go-login-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAdminHandler hanya bisa diakses oleh admin
func RegisterAdminHandler(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pastikan admin hanya bisa membuat akun admin lain
	req.Role = "admin"

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Simpan admin baru ke database
	adminUser := helper.ToUser(req, string(hashedPassword))

	if err := config.DB.Create(&adminUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	// Gunakan helper untuk membuat response
	response := helper.ToUserResponse(adminUser)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin registered successfully!",
		"user":    response,
	})
}
