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

// Types to use for binding (entities with Json Tags)
type BlogDTO struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Owner   Domain.User `json:"owner"`
	Tags    string      `json:"tags"`
	Date    string      `json:"date"`
}

func NewBlogController(Uc usecases.BlogUseCaseI) *BlogController {
	return &BlogController{
		UseCase: Uc,
	}
}

func (BlgCtrl *BlogController) CreateBlogController(c *gin.Context) {
	var blog BlogDTO
	err := c.ShouldBind(&blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// validate if the user is authorized and authenticated
	err = BlgCtrl.UseCase.CreateBlogUC(BlgCtrl.ChangeToDomain(blog))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "blog created successfully"})
}

func (BlgCtrl *BlogController) UpdateBlogController(c *gin.Context) {
	var updated_blog BlogDTO
	err := c.ShouldBindJSON(&updated_blog)

	// Handle binding errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Call usecase and handle different errors
	err = BlgCtrl.UseCase.UpdateBlogUC(BlgCtrl.ChangeToDomain(updated_blog))
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

func (BlgCtrl *BlogController) ChangeToDomain(b BlogDTO) Domain.Blog {
	blog := Domain.Blog{
		ID:      b.ID,
		Date:    b.Date,
		Title:   b.Title,
		Owner:   b.Owner,
		Content: b.Content,
		Tags:    b.Tags,
	}
	return blog
}
