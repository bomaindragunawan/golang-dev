package middleware

import (
	"fmt"
	"go-login-api/config"
	"go-login-api/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware untuk validasi token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required or malformed"})
			c.Abort()
			return
		}

		// Hilangkan "Bearer " agar hanya mendapatkan token asli
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Cek apakah token ada di blacklist
		var count int64
		config.DB.Model(&models.BlacklistToken{}).Where("token = ?", tokenString).Count(&count)
		if count > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// Ambil secret key dari environment
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			fmt.Println("JWT_SECRET is missing!")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Server error: JWT_SECRET is not set"})
			c.Abort()
			return
		}

		// Parse token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Ambil claims dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Ambil user_id dan email dari token
		userID, exists := claims["user_id"]
		email, emailExists := claims["email"]

		if !exists || !emailExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
			c.Abort()
			return
		}

		// Konversi userID ke format yang benar (float64 → uint)
		userIDUint := uint(userID.(float64))

		// Ambil data user dari database
		var user models.User
		if err := config.DB.First(&user, userIDUint).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Simpan user_id dan role ke context agar bisa diakses di middleware RBAC
		c.Set("user_id", userIDUint)
		c.Set("user_role", user.Role)
		c.Set("email", email)

		c.Next()
	}
}
