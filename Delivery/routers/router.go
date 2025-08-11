package routers

import (
	"blog_api/Delivery/controllers"
	infrastructure "blog_api/Infrastructure"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func SetupRouter(BlogCtrl *controllers.BlogController, UserCtrl *controllers.UserController, middleware *infrastructure.AuthMiddleware) {
	// Initialize a new router
	router := gin.Default()

	// load environment
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file failed to load")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")

	if clientID == "" || clientCallbackURL == "" || clientSecret == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required.")
	}

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)

	// Set endpoints
	blogRoutes := router.Group("/blog")
	{
		blogRoutes.GET("/", BlogCtrl.GetAllBlogController)
		blogRoutes.GET("/search", BlogCtrl.SearchBlogController)
		blogRoutes.GET("/filter", BlogCtrl.FilterBlogController)
		blogRoutes.GET("/:id", BlogCtrl.GetBlogController)
		blogRoutes.GET("/:id/view", BlogCtrl.ViewBlogController)
		blogRoutes.GET("/:id/likes", BlogCtrl.LikesController)
		blogRoutes.GET("/:id/dislikes", BlogCtrl.DislikesController)
		blogRoutes.GET("/popular", BlogCtrl.GetPopularBlogs)

		// Authenticated Routes
		authBlog := blogRoutes.Group("/")
		authBlog.Use(middleware.Auth_token())
		{
			authBlog.POST("/", BlogCtrl.CreateBlogController)
			authBlog.PUT("/", BlogCtrl.UpdateBlogController)
			authBlog.DELETE("/:id", BlogCtrl.DeleteBlogController)
			authBlog.GET("/:id/like", BlogCtrl.LikeBlogController)
			authBlog.GET("/:id/dislike", BlogCtrl.DisLikeBlogController)
			authBlog.POST("/:id/comments", BlogCtrl.CommentsBlogController)
			authBlog.POST("/chat", BlogCtrl.AiChatBlogController)
			authBlog.GET("/read_later", BlogCtrl.ReadLatersBlogController)
			authBlog.POST("/:id/read_later", BlogCtrl.InsertReadLatersBlogController)
		}
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/", UserCtrl.RegisterController)
		userRoutes.POST("/verify-otp", UserCtrl.VerifyOTPController)
		userRoutes.POST("/login", UserCtrl.LoginController)
		userRoutes.POST("/forgot-password", UserCtrl.ForgotPasswordController)
		userRoutes.POST("/reset-password", UserCtrl.ResetPasswordController)
		userRoutes.GET("/auth/:provider", UserCtrl.SignInWithProvider)
		userRoutes.GET("/auth/:provider/callback", UserCtrl.OauthCallback)
		userRoutes.POST("/refresh", UserCtrl.RefreshController)

		// Authenticated Routes
		authUser := userRoutes.Group("/")
		authUser.Use(middleware.Auth_token())
		{
			authUser.PUT("/", UserCtrl.UpdateProfileController)
			authUser.POST("/logout", UserCtrl.LogoutController)

			// Admin Routes
			authUser.PUT("/role", middleware.Require_Admin(), UserCtrl.UpdateUserRoleController)
		}
	}
	// Run the router
	router.Run()
}
