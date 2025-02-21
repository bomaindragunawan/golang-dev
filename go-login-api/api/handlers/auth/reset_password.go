package auth

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Handler untuk Reset Password
func ResetPasswordHandler(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		OTP         string `json:"otp" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Cari OTP di database
	var passwordReset models.PasswordReset
	result := config.DB.Where("email = ? AND otp = ?", req.Email, req.OTP).First(&passwordReset)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Cek apakah OTP sudah expired
	if time.Now().After(passwordReset.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
		return
	}

	// Hash password baru
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	// Update password user di database
	config.DB.Model(&models.User{}).Where("email = ?", req.Email).Update("password", string(hashedPassword))

	// Hapus OTP setelah digunakan
	config.DB.Delete(&passwordReset)

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully reset"})
}
