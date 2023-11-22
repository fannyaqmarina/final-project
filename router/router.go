package router

import (
	"final-assignment/controllers"
	"final-assignment/middleware"

	"github.com/gin-gonic/gin"
)

func RouterApp() *gin.Engine {
	router := gin.Default()

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", controllers.Signup)
		authRouter.POST("/login", controllers.Login)
	}

	productRouter := router.Group("/products")
	{
		productRouter.GET("/", controllers.GetAllProduct)
		productRouter.GET("/:productId", controllers.GetOneProduct)

		variantRouter := productRouter.Group("/variants")
		{
			variantRouter.GET("/", controllers.GetAllVariant)
			variantRouter.GET("/:variantId", controllers.GetOneVariant)

			variantRouter.Use(middleware.Authentication())
			variantRouter.POST("/", controllers.CreateVariant)
		}
		productRouter.Use(middleware.Authentication())
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.PUT("/:productId", middleware.ProductAuthorization(), controllers.UpdateProduct)
	}

	return router
}
