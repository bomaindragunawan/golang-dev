package auth

import (
	"fmt"
	"go-login-api/config"
	"go-login-api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Struktur untuk response login
type LoginResponse struct {
	Token string `json:"token"`
}

// Fungsi untuk generate JWT token
func generateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 2).Unix(), // Token berlaku 2 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		fmt.Println("JWT_SECRET is missing!")
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	signedToken, err := token.SignedString([]byte(secretKey))

	// Debugging: Print token yang dibuat
	//fmt.Println("Generated Token:", signedToken)

	return signedToken, err
}

// Handler Login
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//debugging
	//fmt.Println("Email Input:", req.Email)
	//fmt.Println("Password Input:", req.Password)

	// Cek apakah user ada di database
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		fmt.Println("User Not Found in Database!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Debugging: Cetak password yang tersimpan di database
	//fmt.Println("Password Input:", req.Password)
	//fmt.Println("Password Hash (DB):", user.Password)

	// Verifikasi password dengan bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println("Password Mismatch!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Buat token JWT
	accessToken, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Buat refresh token
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Simpan refresh token ke database
	user.RefreshToken = refreshToken
	config.DB.Save(&user)

	// Kirim token ke user
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

func generateRefreshToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Refresh token berlaku 7 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Cari user berdasarkan refresh token
	var user models.User
	result := config.DB.Where("refresh_token = ?", req.RefreshToken).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Buat access token baru
	newAccessToken, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
