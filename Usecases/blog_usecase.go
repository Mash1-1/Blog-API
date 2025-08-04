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
	GetAllBlogUC(limit int, offset int) ([]Domain.Blog, error)
	SearchBlogUC(Domain.Blog) ([]Domain.Blog, error)
	DeleteBlogUC(string) error
	FilterBlogUC(Domain.Blog) ([]Domain.Blog, error)
	GetByIdBlogUC(string) (Domain.Blog, error)
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
	var tmp = Domain.User{}
	if searchBlog.Title == "" && searchBlog.Owner == tmp {
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
