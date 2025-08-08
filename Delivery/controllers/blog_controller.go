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

// Types to use for binding (entities with Json Tags) and also bson format for storing
type BlogDTO struct {
	ID          string    `json:"id" bson:"ID"`
	Title       string    `json:"title" bson:"Title"`
	Content     string    `json:"content" bson:"Content"`
	Owner_email string    `json:"owner" bson:"Owner"`
	Tags        []string  `json:"tags" bson:"Tags"`
	Date        time.Time `json:"date" bson:"Date"`
	Likes       int       `json:"likes" bson:"Likes"`
	Dislikes    int       `json:"dislikes" bson:"Dislikes"`
	ViewCount   int       `json:"viewCount" bson:"ViewCount"`
	Comments    []string  `json:"comments" bson:"Comments"`
}

func NewBlogController(Uc usecases.BlogUseCaseI) *BlogController {
	return &BlogController{
		UseCase: Uc,
	}
}

func (BlgCtrl *BlogController) CreateBlogController(c *gin.Context) {
	var blog BlogDTO
	err := c.ShouldBindJSON(&blog)
	user, _ := c.Get("user")
	blog.Owner_email = user.(*Domain.User).Email
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
	c.JSON(http.StatusCreated, gin.H{"message: ": "blog created successfully"})
}

func (BlgCtrl *BlogController) SearchBlogController(c *gin.Context) {
	var SearchBlog BlogDTO
	err := c.ShouldBindJSON(&SearchBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Validate request using usecase
	blogs, err := BlgCtrl.UseCase.SearchBlogUC(BlgCtrl.ChangeToDomain(SearchBlog))
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
	var updated_blog BlogDTO
	err := c.ShouldBindJSON(&updated_blog)

	// Handle binding errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Call usecase and handle different errors
	err = BlgCtrl.UseCase.UpdateBlogUC(BlgCtrl.ChangeToDomain(updated_blog))
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
	var FilterBlog BlogDTO
	err := c.ShouldBindJSON(&FilterBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if FilterBlog.Date.IsZero() && FilterBlog.Tags == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both Date and Tags can't be empty"})
		return
	}
	blogs, err := BlgCtrl.UseCase.FilterBlogUC(BlgCtrl.ChangeToDomain(FilterBlog))
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

// method to convert from Blog DTO to Blog structure
func (BlgCtrl *BlogController) ChangeToDomain(BlgDto BlogDTO) Domain.Blog {
	blog := Domain.Blog{
		ID:          BlgDto.ID,
		Date:        BlgDto.Date,
		Title:       BlgDto.Title,
		Owner_email: BlgDto.Owner_email,
		Content:     BlgDto.Content,
		Tags:        BlgDto.Tags,
		Likes:       BlgDto.Likes,
		Dislikes:    BlgDto.Dislikes,
		ViewCount:   BlgDto.ViewCount,
		Comments:    BlgDto.Comments,
	}
	return blog
}
