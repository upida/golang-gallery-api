package router

import (
	"gallery/controllers"
	"gallery/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	public := router.Group("")
	{
		user := public.Group("/users")
		{
			user.POST("/register", controllers.Register)
			user.POST("/login", controllers.Login)
		}

		photo := public.Group("/photos")
		{
			photo.GET("/:uuid", controllers.GetPhoto)
		}
	}

	protected := router.Group("")
	protected.Use(middlewares.JwtAuthMiddleware())
	{
		user := protected.Group("/users")
		{
			user.GET("", controllers.CurrentUser)
			user.PUT("", controllers.UpdateUser)
		}

		photo := protected.Group("/photos")
		{
			photo.GET("", controllers.UserPhotos)
			photo.POST("uploads", controllers.UploadPhotos)
			photo.DELETE("/:uuid", controllers.DeletePhoto)
		}
	}

}
