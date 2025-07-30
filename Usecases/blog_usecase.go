package usecases

import (
	"blog_api/Domain"
	"blog_api/Repositories"
)

type BlogUseCase struct {
	Repository Repositories.BlogRepositoryI
}
type BlogUseCaseI interface {
	CreateBlog(Domain.Blog) error
}

func NewBlogUseCase(Repo Repositories.BlogRepositoryI) *BlogUseCase {
	return &BlogUseCase{
		Repository: Repo,
	}
}

func (BlgUseCase *BlogUseCase) CreateBlog(blog Domain.Blog) error {
	err := BlgUseCase.Repository.Create(blog)
	return err
}
