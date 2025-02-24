package auth

import (
	"go-login-api/config"
	"go-login-api/helper"
	"go-login-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest // Pakai struct dari models

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Default role jika tidak ditentukan
	if req.Role == "" {
		req.Role = "user"
	}

	// **Batasi hanya admin yang bisa membuat admin baru**
	if req.Role == "admin" {
		_, isAdmin := c.Get("user_role") // Cek role dari middleware JWT
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create an admin user"})
			return
		}
	}

	// Gunakan helper untuk konversi RegisterRequest ke User
	user := helper.ToUser(req, string(hashedPassword))

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Gunakan struct `UserResponse` untuk response
	// **Gunakan helper untuk konversi User ke UserResponse**
	response := helper.ToUserResponse(user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully!",
		"user":    response,
	})
}
