package usecases

import (
	"blog_api/Domain"
	"errors"
)

type BlogUseCase struct {
	Repository Domain.BlogRepositoryI
}
type BlogUseCaseI interface {
	CreateBlogUC(Domain.Blog) error
	UpdateBlogUC(Domain.Blog) error
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

func (BlgUC *BlogUseCase) UpdateBlogUC(updatedBlog Domain.Blog) error {
	// Handle empty blog update
	if updatedBlog.Content == "" && updatedBlog.Title == "" && updatedBlog.Tags == nil {
		return errors.New("can't update into empty blog")
	}
	return BlgUC.Repository.UpdateBlog(&updatedBlog)
}
