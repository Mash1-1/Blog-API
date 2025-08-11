package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog_api/Domain"
)

type BlogController struct {
	UseCase Domain.BlogUseCaseI
}

func NewBlogController(Uc Domain.BlogUseCaseI) *BlogController {
	return &BlogController{
		UseCase: Uc,
	}
}

func (BlgCtrl *BlogController) CreateBlogController(c *gin.Context) {
	var blog BlogDTO
	log.Print("gets here")
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
	_, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	// check if the user have liked the post previously
	user, _ := c.Get("user")
	liked, err := BlgCtrl.UseCase.CheckIfLiked(user.(*Domain.User).Email, id)
	if err != nil {
		if err.Error() == "invalid blog id or user email when checking liked" {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	var Liketrk Domain.LikeTracker
	Liketrk.BlogID = id
	Liketrk.UserEmail = user.(*Domain.User).Email

	if liked == 1 {
		Liketrk.Liked = 0
		BlgCtrl.UseCase.AddLikeUC(Liketrk)
		c.JSON(http.StatusOK, gin.H{"message": "removed like from this blog"})
		return
	}

	// Add the like tracked to the database
	Liketrk.Liked = 1
	err = BlgCtrl.UseCase.AddLikeUC(Liketrk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message: ": "liked"})
}

func (BlgCtrl *BlogController) DisLikeBlogController(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("user")
	user_email := user.(*Domain.User).Email
	_, err := BlgCtrl.UseCase.GetByIdBlogUC(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	// Check if user has already disliked this post
	var Liketrk Domain.LikeTracker
	Liketrk.BlogID = id
	Liketrk.UserEmail = user_email
	liked, err := BlgCtrl.UseCase.CheckIfLiked(user_email, id)
	if err != nil {
		if err.Error() == "invalid blog id or user email when checking liked" {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if liked == -1 {
		Liketrk.Liked = 0
		err := BlgCtrl.UseCase.AddLikeUC(Liketrk)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "removed dislike from blog"})
		return
	}
	Liketrk.Liked = -1
	err = BlgCtrl.UseCase.AddLikeUC(Liketrk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message: ": "Disliked"})
}

func (BlgCtrl *BlogController) LikesController(c *gin.Context) {
	id := c.Param("id")
	num, err := BlgCtrl.UseCase.Likes(id)
	if err != nil {
		if err.Error() == "id field can not be empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"number of likes": num})
}

func (BlgCtrl *BlogController) DislikesController(c *gin.Context) {
	id := c.Param("id")
	num, err := BlgCtrl.UseCase.Dislikes(id)
	if err != nil {
		if err.Error() == "id field can not be empty" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"number of dislikes": num})
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
	log.Print("gets here")
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

func (BlgCtrl *BlogController) GetPopularBlogs(c *gin.Context) {
	blogs, err := BlgCtrl.UseCase.GetPopularBlogs()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blogs: ": blogs})
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
		ViewCount: BlgDto.ViewCount,
		Comments:  BlgDto.Comments,
	}
	return blog
}