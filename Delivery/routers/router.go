package routers

import (
	"blog_api/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(BlogCtrl *controllers.BlogController) {
	// Initialize a new router
	router := gin.Default()

	// Set endpoints
	router.PUT("/update_blog", BlogCtrl.UpdateBlogController)
	router.POST("/create_blog", BlogCtrl.CreateBlogController)
	router.GET("/search_blog", BlogCtrl.SearchBlogController)
	router.DELETE("/delete_blog/:id", BlogCtrl.DeleteBlogController)

	// Run the router
	router.Run()
}
