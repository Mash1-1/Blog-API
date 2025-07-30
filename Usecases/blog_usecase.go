package usecases

import (
	"blog_api/Domain"
	"blog_api/Repositories"
	"errors"
)

type BlogUseCase struct {
	Repository Repositories.BlogRepositoryI
}
type BlogUseCaseI interface {
	UpdateBlogUC(Domain.Blog) error
}

func NewBlogUseCase(Repo Repositories.BlogRepositoryI) *BlogUseCase {
	return &BlogUseCase{
		Repository: Repo,
	}
}

func (BlgUC *BlogUseCase) UpdateBlogUC(updatedBlog Domain.Blog) error {
	tmp := Domain.Blog{}
	// Handle empty blog update
	if updatedBlog == tmp {
		return errors.New("can't update into empty blog")
	}
	return BlgUC.Repository.UpdateBlog(&updatedBlog)
}