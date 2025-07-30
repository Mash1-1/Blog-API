package controllers

import (
	"blog_api/Domain"
	usecases "blog_api/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	UseCase usecases.BlogUseCaseI
}

func NewBlogController(Uc usecases.BlogUseCaseI) *BlogController {
	return &BlogController{
		UseCase: Uc,
	}
}

func (BlgCtrl *BlogController) UpdateBlogController(c *gin.Context) {}

func (BlgCtrl *BlogController) CreateBlogController(c *gin.Context) {
	// confirm the user is authenticated and authorized
	var blog Domain.Blog
	err := c.BindJSON(&blog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Incorrect Request Body."})
		return
	}
	// if all passed insert a new blog to the database
	err = BlgCtrl.UseCase.CreateBlog(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	//return success message
	c.JSON(http.StatusOK, gin.H{"Message: ": "Blog created succesfully."})
}
