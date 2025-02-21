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

		// Debugging: Print token yang diterima
		//fmt.Println("Token Diterima:", tokenString)

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			fmt.Println("Token tidak ada atau format salah")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required or malformed"})
			c.Abort()
			return
		}

		// Hilangkan "Bearer " agar hanya mendapatkan token asli
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Debugging: Print token setelah pembersihan
		//fmt.Println("Token After Cleaning:", tokenString)

		// Cek apakah token ada di blacklist
		var count int64
		config.DB.Model(&models.BlacklistToken{}).Where("token = ?", tokenString).Count(&count)
		if count > 0 {
			//fmt.Println("Token masuk blacklist, akses ditolak")
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

		// Debug: Print secret key
		//fmt.Println("JWT Secret:", secretKey)

		// Parse token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		// Debug: Print error parsing token
		if err != nil {
			fmt.Println("JWT Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Validasi token
		if !token.Valid {
			fmt.Println("Token tidak valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Ambil claims dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Claims token tidak valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Debugging: Print claims dari token
		//fmt.Println("Claims:", claims)

		// Simpan email user ke context Gin
		c.Set("email", claims["email"])
		c.Next()
	}
}
