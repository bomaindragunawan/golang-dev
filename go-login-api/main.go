package main

import (
	"go-login-api/config"
	"go-login-api/models"
	"go-login-api/routes"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})
	r := routes.Setup()
	r.Run(":8080")
}
