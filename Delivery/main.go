package main

import (
	"blog_api/Delivery/controllers"
	"blog_api/Delivery/routers"
	infrastructure "blog_api/Infrastructure"
	"blog_api/Repositories"
	usecases "blog_api/Usecases"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize controllers and router
	db := Repositories.InitializeDb()

	// blog dependency injection
	blog_repo := Repositories.NewBlogRepository(db)
	blog_usecase := usecases.NewBlogUseCase(blog_repo)
	blog_controller := controllers.NewBlogController(blog_usecase)

	// Get required email info from the env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load environment variables")
	}
	Host := os.Getenv("SMTP_HOST")
	Port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	Username := os.Getenv("SMTP_USERNAME")
	Pass := os.Getenv("SMTP_PASSWORD")
	frm := os.Getenv("SMTP_FROM")

	j_serv := infrastructure.Jwt_serv{}
	generator_otp := infrastructure.Generator{}
	password_service := infrastructure.PasswordService{}
	mailr := infrastructure.NewMailer(Host, Port, Username, Pass, frm)

	// user dependency injection
	user_repo := Repositories.NewUserRepository(db)
	user_usecase := usecases.NewUserUsecase(user_repo, password_service, &mailr, generator_otp, j_serv)

	// auth middleware
	middleware := infrastructure.AuthMiddleware{Usecase: user_usecase}
	user_controller := controllers.NewUserController(user_usecase)

	// router
	routers.SetupRouter(blog_controller, &user_controller, &middleware)
}
