package usecases

import "blog_api/Repositories"

type BlogUseCase struct {
	Repository Repositories.BlogRepositoryI
}
type BlogUseCaseI interface {}

func NewBlogUseCase(Repo Repositories.BlogRepositoryI) *BlogUseCase {
	return &BlogUseCase{
		Repository: Repo,
	}
}