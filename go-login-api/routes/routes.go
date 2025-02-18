package routes

import (
	"go-login-api/handlers"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)
	return router
}
