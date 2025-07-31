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

// Types to use for binding (entities with Json Tags) and also bson format for storing
type BlogDTO struct {
	ID      string    `json:"id" bson:"ID"`
	Title   string    `json:"title" bson:"Title"`
	Content string    `json:"content" bson:"Content"`
	Owner   User      `json:"owner" bson:"Owner"`
	Tags    []string  `json:"tags" bson:"Tags"`
	Date    time.Time `json:"date" bson:"Date"`
}

// method to convert from Blog DTO to Blog structure
func (BlgDto *BlogDTO) ToDomain() Blog {
	blog := Blog{
		ID:      BlgDto.ID,
		Date:    BlgDto.Date,
		Title:   BlgDto.Title,
		Owner:   BlgDto.Owner,
		Content: BlgDto.Content,
		Tags:    BlgDto.Tags,
	}
	return blog
}

// method to convert Blog struct to BlogDTO object
func (Blg *Blog) ToBlogDTO() BlogDTO {
	blogDTO := BlogDTO{
		ID:      Blg.ID,
		Title:   Blg.Title,
		Content: Blg.Content,
		Owner:   Blg.Owner,
		Tags:    Blg.Tags,
		Date:    Blg.Date,
	}
	return blogDTO
}

type BlogRepositoryI interface {
	Create(blog *Blog) error
	UpdateBlog(updatedBlog *Blog) error
	SearchBlog(searchBlog *Blog) ([]Blog, error)
	DeleteBlog(id string) error
}
