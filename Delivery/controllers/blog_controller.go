package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"blog_api/Domain"
	usecases "blog_api/Usecases"
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
	Tags    []string    `json:"tags"`
	Date    time.Time   `json:"date"`
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

func (BlgCtrl *BlogController) GetAllBlogController(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
		return
	}

	blogs, err := BlgCtrl.UseCase.GetAllBlogUC(int(limit), int(offset))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}
