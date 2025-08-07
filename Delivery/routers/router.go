package routers

import (
	"blog_api/Delivery/controllers"
	infrastructure "blog_api/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(BlogCtrl *controllers.BlogController, UserCtrl *controllers.UserController, middleware *infrastructure.AuthMiddleware) {
	// Initialize a new router
	router := gin.Default()

	// Set endpoints
	blogRoutes := router.Group("/blog")
	{
		blogRoutes.GET("/", BlogCtrl.GetAllBlogController)
		blogRoutes.POST("/", middleware.Auth_token(), BlogCtrl.CreateBlogController)
		blogRoutes.PUT("/", middleware.Auth_token(), BlogCtrl.UpdateBlogController)
		blogRoutes.DELETE("/:id", middleware.Auth_token(), BlogCtrl.DeleteBlogController)
		blogRoutes.GET("/search", BlogCtrl.SearchBlogController)
		blogRoutes.GET("/filter", BlogCtrl.FilterBlogController)
		blogRoutes.GET("/:id", BlogCtrl.GetBlogController)
		blogRoutes.GET("/:id/likes", middleware.Auth_token(), BlogCtrl.LikeBlogController)
		blogRoutes.GET("/:id/dislikes", middleware.Auth_token(), BlogCtrl.DisLikeBlogController)
		blogRoutes.GET("/:id/view", BlogCtrl.ViewBlogController)
		blogRoutes.POST("/:id/comments", middleware.Auth_token(), BlogCtrl.CommentsBlogController)
		blogRoutes.POST("/chat", middleware.Auth_token(), BlogCtrl.AiChatBlogController)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/", UserCtrl.RegisterController)
		userRoutes.POST("/verify-otp", UserCtrl.VerifyOTPController)
		userRoutes.POST("/login", UserCtrl.LoginController)
		userRoutes.POST("/forgot-password", UserCtrl.ForgotPasswordController)
		userRoutes.POST("/reset-password", UserCtrl.ResetPasswordController)
	}
	// Run the router
	router.Run()
}
