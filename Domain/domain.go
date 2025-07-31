package Domain

import "time"

type User struct {
	Username string
	Email    string
	Password string
	Bio      string
	Role     string
}

type Blog struct {
	ID      string
	Title   string
	Content string
	Owner   User
	Tags    []string
	Date    time.Time
}

type BlogRepositoryI interface {
	Create(blog *Blog) error
	UpdateBlog(updatedBlog *Blog) error
	SearchBlog(searchBlog *Blog) ([]Blog, error)
}
