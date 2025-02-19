package handlers

import (
	"go-login-api/config"
	"go-login-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debugging: Cetak password sebelum hashing
	//fmt.Println("Raw Password:", user.Password)

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Debugging: Cetak password setelah hashing
	//fmt.Println("Hashed Password:", string(hashedPassword))

	user.Password = string(hashedPassword)

	config.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})

}
