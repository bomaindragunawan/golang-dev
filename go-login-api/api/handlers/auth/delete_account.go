package auth

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteAccountHandler memungkinkan user untuk menghapus akunnya sendiri
func DeleteAccountHandler(c *gin.Context) {
	// Ambil user ID dari context (di-set oleh middleware JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Hapus user dari database
	if err := config.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// DeleteUserByAdminHandler memungkinkan admin menghapus akun user lain
func DeleteUserByAdminHandler(c *gin.Context) {
	// Ambil user ID dari parameter URL
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Cek apakah user dengan ID tersebut ada
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hapus user dari database
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
