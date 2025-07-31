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

func (BlgCtrl *BlogController) CreateBlogController(c *gin.Context) {
	var blog Domain.BlogDTO
	err := c.ShouldBind(&blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// validate if the user is authorized and authenticated
	err = BlgCtrl.UseCase.CreateBlogUC(blog.ToDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "blog created successfully"})
}

func (BlgCtrl *BlogController) SearchBlogController(c *gin.Context) {
	var SearchBlog Domain.BlogDTO
	err := c.ShouldBindJSON(&SearchBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Validate request using usecase
	blogs, err := BlgCtrl.UseCase.SearchBlogUC(SearchBlog.ToDomain())
	if err != nil {
		if err.Error() == "can't update into empty blog" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Blogs found": blogs})
}

func (BlgCtrl *BlogController) UpdateBlogController(c *gin.Context) {
	var updated_blog Domain.BlogDTO
	err := c.ShouldBindJSON(&updated_blog)

	// Handle binding errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Call usecase and handle different errors
	err = BlgCtrl.UseCase.UpdateBlogUC(updated_blog.ToDomain())
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "blog updated successfuly"})
		return
	}
	if err.Error() == "blog not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err.Error() == "can't update into empty blog" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (BlgCtrl *BlogController) DeleteBlogController(c *gin.Context) {
	id := c.Param("id")

	//validate if the user is owner or admin
	if err := BlgCtrl.UseCase.DeleteBlogUC(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": " blog deleted successfully"})
}
