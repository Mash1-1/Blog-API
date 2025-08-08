package usecases

import (
	"blog_api/Domain"
	infrastructure "blog_api/Infrastructure"
	"errors"
	"strings"
)

type BlogUseCase struct {
	Repository Domain.BlogRepositoryI
}
type BlogUseCaseI interface {
	CreateBlogUC(Domain.Blog) error
	UpdateBlogUC(Domain.Blog) error
	GetAllBlogUC(limit int, offset int) ([]Domain.Blog, error)
	SearchBlogUC(Domain.Blog) ([]Domain.Blog, error)
	DeleteBlogUC(string) error
	FilterBlogUC(Domain.Blog) ([]Domain.Blog, error)
	GetByIdBlogUC(string) (Domain.Blog, error)
	AIChatBlogUC(Domain.ChatRequest) (*string, error)
}

func NewBlogUseCase(Repo Domain.BlogRepositoryI) *BlogUseCase {
	return &BlogUseCase{
		Repository: Repo,
	}
}

func (BlgUseCase *BlogUseCase) CreateBlogUC(blog Domain.Blog) error {
	err := BlgUseCase.Repository.Create(&blog)
	return err
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
