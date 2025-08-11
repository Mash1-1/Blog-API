package usecases

import (
	"blog_api/Domain"
	infrastructure "blog_api/Infrastructure"
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type BlogUseCase struct {
	Repository Domain.BlogRepositoryI
}

func NewBlogUseCase(Repo Domain.BlogRepositoryI) *BlogUseCase {
	return &BlogUseCase{
		Repository: Repo,
	}
}

func (BlgUseCase *BlogUseCase) CreateBlogUC(blog Domain.Blog) error {
	blog.ID = uuid.New().String()
	err := BlgUseCase.Repository.Create(&blog)
	return err
}

func (BlgUseCase *BlogUseCase) AddLikeUC(lt Domain.LikeTracker) error {
	// Delete previous instances
	err := BlgUseCase.Repository.DeleteLikeTk(lt)
	if err != nil {
		return err
	}
	return BlgUseCase.Repository.CreateLikeTk(lt)
}

func (BlgUseCase *BlogUseCase) CheckIfLiked(user_email, blogId string) (int, error) {
	if user_email == "" || blogId == "" {
		return 0, errors.New("invalid blog id or user email when checking liked")
	}
	liked, err := BlgUseCase.Repository.FindLiked(user_email, blogId)
	if err != nil && err.Error() == "mongo: no documents in result" {
		// If user hasnt liked this post before, create a new doc to like it
		var Liketrk Domain.LikeTracker
		Liketrk.BlogID = blogId
		Liketrk.UserEmail = user_email
		Liketrk.Liked = 0

		err := BlgUseCase.Repository.CreateLikeTk(Liketrk)
		if err != nil {
			return 0, err
		}
		return 0, nil
	} else {
		if err != nil {
			return 0, err
		}
	}
	return liked.Liked, err
}

func (BlgUsecase *BlogUseCase) Likes(id string) (int64, error) {
	if id == "" {
		return 0, errors.New("id field can not be empty")
	}
	_, err := BlgUsecase.Repository.GetBlog(id)
	if err != nil {
		return 0, err
	}

	return BlgUsecase.Repository.NumberOfLikes(id)
}

func (BlgUsecase *BlogUseCase) Dislikes(id string) (int64, error) {
	if id == "" {
		return 0, errors.New("id field can not be empty")
	}
	_, err := BlgUsecase.Repository.GetBlog(id)
	if err != nil {
		return 0, err
	}
	return BlgUsecase.Repository.NumberOfDislikes(id)
}

func (BlgUseCase *BlogUseCase) SearchBlogUC(searchBlog Domain.Blog) ([]Domain.Blog, error) {
	// Check if required fields are available
	if searchBlog.Title == "" && searchBlog.Owner_email == "" {
		return []Domain.Blog{}, errors.New("can't search for blog with empty searching fileds.(Title or Owner)")
	}
	return BlgUseCase.Repository.SearchBlog(&searchBlog)
}

func (BlgUC *BlogUseCase) UpdateBlogUC(updatedBlog Domain.Blog) error {
	// Handle empty blog update
	if updatedBlog.Content == "" && updatedBlog.Title == "" && updatedBlog.Tags == nil {
		return errors.New("can't update into empty blog")
	}
	return BlgUC.Repository.UpdateBlog(&updatedBlog)
}

func (BlgUseCase *BlogUseCase) GetAllBlogUC(limit int, offset int) ([]Domain.Blog, error) {
	return BlgUseCase.Repository.GetAllBlogs(limit, offset)
}

func (BlgUC *BlogUseCase) DeleteBlogUC(id string) error {
	err := BlgUC.Repository.DeleteBlog(id)
	return err
}

func (BlgUseCase *BlogUseCase) FilterBlogUC(filterBlog Domain.Blog) ([]Domain.Blog, error) {
	return BlgUseCase.Repository.FilterBlog(&filterBlog)
}

func (BlgUseCase *BlogUseCase) GetByIdBlogUC(id string) (Domain.Blog, error) {
	return BlgUseCase.Repository.GetBlog(id)
}

func (BlgUseCase *BlogUseCase) AIChatBlogUC(message Domain.ChatRequest) (*string, error) {
	AI := infrastructure.InitAIClient()

	blog_text, err := AI.Generate_blog_content(message.Message)

	if err != nil {
		return nil, err
	}

	return blog_text, nil
}

func (BlgUseCase *BlogUseCase) GetPopularBlogs() ([]Domain.Blog, error) {
	// Get all the blogs and sort them by their popularity score
	blogs, err := BlgUseCase.Repository.GetAllBlogs(0, 0)
	if err != nil {
		return []Domain.Blog{}, err
	}
	sort.Slice(blogs, func(i, j int) bool {

	// Calculate popularity scores
	likesI, _ := BlgUseCase.Repository.NumberOfLikes(blogs[i].ID)
	dislikesI, _ := BlgUseCase.Repository.NumberOfDislikes(blogs[i].ID)
	scoreI := likesI + int64(blogs[i].ViewCount) - dislikesI

	likesJ, _ := BlgUseCase.Repository.NumberOfLikes(blogs[j].ID)
	dislikesJ, _ := BlgUseCase.Repository.NumberOfDislikes(blogs[j].ID)
	scoreJ := likesJ + int64(blogs[j].ViewCount) - dislikesJ

	return scoreI > scoreJ // Descending order
})	
	return blogs, nil
	// return BlgUseCase.Repository.GetPopularBlogs()
}

func RemoveLinesContaining(text string) string {
	phrases := []string{"Okay, here's", " I'll try", "Feel free to give me", "Let me know what you think", "The more information you give me, the better I can tailor", "?", "**", "I hope this helps", "Let me know if you'd like me to create", "("}

	lines := strings.Split(text, "\n")
	cleanedLines := []string{}

	for _, line := range lines {
		shouldSkip := false
		lowerLine := strings.ToLower(line)

		for _, phrase := range phrases {
			if strings.Contains(lowerLine, strings.ToLower(phrase)) {
				shouldSkip = true
				break
			}
		}

		if !shouldSkip && strings.TrimSpace(line) != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "")
}
