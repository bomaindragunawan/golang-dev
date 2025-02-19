package main

import (
	"go-login-api/api/routes"
	"go-login-api/config"
)

func main() {
	config.ConnectDB()
	//config.DB.AutoMigrate(&models.User{})
	r := routes.Setup()
	r.Run(":8080")
}
