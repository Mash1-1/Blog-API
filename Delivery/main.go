package main

import (
	"blog_api/Delivery/controllers"
	"blog_api/Delivery/routers"
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

	routers.SetupRouter(blog_controller)
}