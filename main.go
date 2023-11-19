package main

import (
	"final-assignment/controllers"
	"final-assignment/initializers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/signup", controllers.Signup)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
