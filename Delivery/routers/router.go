package routers

import (
	"blog_api/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(BlogCtrl *controllers.BlogController, UserCtrl *controllers.UserController) {
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
		blogRoutes.GET("/filter", BlogCtrl.FilterBlogController)
		blogRoutes.GET("/:id", BlogCtrl.GetBlogController)
		blogRoutes.GET("/:id/likes", BlogCtrl.LikeBlogController)
		blogRoutes.GET("/:id/dislikes", BlogCtrl.DisLikeBlogController)
		blogRoutes.GET("/:id/view", BlogCtrl.ViewBlogController)
		blogRoutes.POST("/:id/comments", BlogCtrl.CommentsBlogController)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/", UserCtrl.RegisterController)
		userRoutes.POST("/verify-otp", UserCtrl.VerifyOTPController)
		userRoutes.POST("/login", UserCtrl.LoginController)
	}
	// Run the router
	router.Run()
}
