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
	var blog Domain.BlogDTO
	err := c.ShouldBindJSON(&blog)

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
	c.JSON(http.StatusCreated, gin.H{"message: ": "blog created successfully"})
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
	if err != nil {
		if err.Error() == "blog not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "can't update into empty blog" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "blog updated successfuly"})
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

func (BlgCtrl *BlogController) FilterBlogController(c *gin.Context) {
	var FilterBlog Domain.BlogDTO
	err := c.ShouldBindJSON(&FilterBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if FilterBlog.Date.IsZero() && FilterBlog.Tags == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both Date and Tags can't be empty"})
		return
	}
	blogs, err := BlgCtrl.UseCase.FilterBlogUC(FilterBlog.ToDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// filter by popularity is not implemented
	c.JSON(http.StatusOK, gin.H{"filtered blogs: ": blogs})
}

func (BlgCtrl *BlogController) GetBlogController(c *gin.Context) {
	id := c.Param("id")
	blog, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Blog: ": blog})
}

func (BlgCtrl *BlogController) LikeBlogController(c *gin.Context) {
	id := c.Param("id")
	blog, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	// check if the user have liked the post previously
	blog.Likes += 1
	err = BlgCtrl.UseCase.UpdateBlogUC(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "liked"})
}

func (BlgCtrl *BlogController) DisLikeBlogController(c *gin.Context) {
	id := c.Param("id")
	blog, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	// check if the user have disliked the post previously
	blog.Dislikes += 1
	err = BlgCtrl.UseCase.UpdateBlogUC(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "Disliked"})
}

func (BlgCtrl *BlogController) ViewBlogController(c *gin.Context) {
	id := c.Param("id")
	blog, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	// check if the user have viewed the post previously
	blog.ViewCount += 1
	err = BlgCtrl.UseCase.UpdateBlogUC(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "View increased"})
}

func (BlgCtrl *BlogController) CommentsBlogController(c *gin.Context) {
	var comment string
	err := c.BindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error: ": err.Error()})
		return
	}
	id := c.Param("id")
	blog, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	blog.Comments = append(blog.Comments, comment)
	err = BlgCtrl.UseCase.UpdateBlogUC(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "comment added"})
}

func (BlgCtrl *BlogController) AiChatBlogController(c *gin.Context) {
	var message Domain.ChatRequest
	err := c.BindJSON(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	var response Domain.ChatResponse
	response.Reply, err = BlgCtrl.UseCase.AIChatBlogUC(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": response})
}
