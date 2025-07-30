package Domain

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
	Tags    string
	Date    string
}