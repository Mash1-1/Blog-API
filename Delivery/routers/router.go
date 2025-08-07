package routers

import (
	"blog_api/Delivery/controllers"
	"log"
	"os"
	infrastructure "blog_api/Infrastructure"

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
		userRoutes.GET("/auth/:provider", UserCtrl.SignInWithProvider)
		userRoutes.GET("/auth/:provider/callback", UserCtrl.OauthCallback)
	}
	// Run the router
	router.Run()
}
