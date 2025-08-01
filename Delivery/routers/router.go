package routers

import (
	"blog_api/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(BlogCtrl *controllers.BlogController) {
	// Initialize a new router
	router := gin.Default()

	// Set endpoints
	blogRoutes := router.Group("/blog")
	{
		blogRoutes.GET("/", BlogCtrl.GetAllBlogController)
		blogRoutes.POST("/", BlogCtrl.CreateBlogController)
		blogRoutes.PUT("/", BlogCtrl.UpdateBlogController)
		blogRoutes.DELETE("/:id", BlogCtrl.DeleteBlogController)
		blogRoutes.GET("/search", BlogCtrl.SearchBlogController)
	}

	// Run the router
	router.Run()
}
