package router

import (
	"final-assignment/controllers"

	"github.com/gin-gonic/gin"
)

func RouterApp() *gin.Engine {
	router := gin.Default()

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", controllers.Signup)
		authRouter.POST("/login", controllers.Login)
	}
	return router
}
