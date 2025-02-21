package auth

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	// Ambil token dari header Authorization
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		//fmt.Println("Logout gagal: Token tidak ditemukan di header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Hapus prefix "Bearer "
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Debugging: Print token yang akan disimpan
	//fmt.Println("Token yang akan dimasukkan ke blacklist:", tokenString)

	// Cek apakah token sudah ada di blacklist
	var count int64
	config.DB.Model(models.BlacklistToken{}).Where("token = ?", tokenString).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Token already blacklisted"})
		return
	}

	// Simpan token ke database (Blacklist Token)
	blacklist := models.BlacklistToken{Token: tokenString}
	saveResult := config.DB.Create(&blacklist)

	// Debugging: Cek apakah token berhasil disimpan
	if saveResult.Error != nil {
		//fmt.Println("Gagal menyimpan token ke blacklist:", saveResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	//fmt.Println("Token berhasil disimpan ke blacklist:", tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
