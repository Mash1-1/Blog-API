package main

import (
	"blog_api/Delivery/controllers"
	"blog_api/Delivery/routers"
	infrastructure "blog_api/Infrastructure"
	"blog_api/Repositories"
	usecases "blog_api/Usecases"
	"fmt"
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

	password_service := infrastructure.PasswordService{}
	user_repo := Repositories.NewUserRepository(user_database)
	user_usecase := usecases.NewUserUsecase(user_repo, password_service)
	user_controller := controllers.NewUserController(user_usecase)

	routers.SetupRouter(blog_controller, &user_controller)
}
