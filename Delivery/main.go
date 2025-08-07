package main

import (
	"blog_api/Delivery/controllers"
	"blog_api/Delivery/routers"
	infrastructure "blog_api/Infrastructure"
	"blog_api/Repositories"
	usecases "blog_api/Usecases"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize controllers and router
	blog_database, err := Repositories.InitializeBlogDB()
	if err != nil {
		fmt.Println("Failed while creating blog database!")
		return
	}
	blog_repo := Repositories.NewBlogRepository(blog_database)
	blog_usecase := usecases.NewBlogUseCase(blog_repo)
	blog_controller := controllers.NewBlogController(blog_usecase)

	user_database, err := Repositories.InitializeUserDB()
	if err != nil {
		fmt.Println("Failed while creating user database!")
		return
	}

	// Get required email info from the env file
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load environment variables")
	}
	Host := os.Getenv("SMTP_HOST")
	Port := os.Getenv("SMTP_PORT")
	Username := os.Getenv("SMTP_USERNAME")
	Pass := os.Getenv("SMTP_PASSWORD")
	frm := os.Getenv("SMTP_FROM")

	j_serv := infrastructure.Jwt_serv{}
	generator_otp := infrastructure.Generator{}
	password_service := infrastructure.PasswordService{}
	mailr := infrastructure.NewMailer(Host, Port, Username, Pass, frm)
	user_repo := Repositories.NewUserRepository(user_database)
	user_usecase := usecases.NewUserUsecase(user_repo, password_service, &mailr, generator_otp, j_serv)
	middleware := infrastructure.AuthMiddleware{Usecase: user_usecase}
	user_controller := controllers.NewUserController(user_usecase)

	routers.SetupRouter(blog_controller, &user_controller, &middleware)
}
