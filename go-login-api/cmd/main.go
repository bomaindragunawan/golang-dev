package main

import (
	"go-login-api/api/routes"
	"go-login-api/config"
	"go-login-api/models"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})
	//config.DB.AutoMigrate(&models.BlacklistToken{})
	//config.DB.AutoMigrate(&models.PasswordReset{})
	r := routes.Setup()
	r.Run(":8080")

}
