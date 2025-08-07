package Domain

import (
	"time"

	"github.com/markbates/goth"
)

type User struct {
	Username string
	Email    string
	Password string
	Bio      string
	Role     string
	Verfied  bool
	OTP      string
	OTPTime  time.Time
	Provider string
}

type Blog struct {
	ID        string
	Title     string
	Content   string
	Owner_email    string
	Tags      []string
	Date      time.Time
	Likes     int
	Dislikes  int
	ViewCount int
	Comments  []string
}


type BlogRepositoryI interface {
	Create(blog *Blog) error
	UpdateBlog(updatedBlog *Blog) error
	GetAllBlogs(limit int, offset int) ([]Blog, error)
	SearchBlog(searchBlog *Blog) ([]Blog, error)
	DeleteBlog(id string) error
	FilterBlog(filterBlog *Blog) ([]Blog, error)
	GetBlog(id string) (Blog, error)
}

type ResetTokenS struct {
	Email       string
	Token       string
	Created_at  time.Time
	NewPassword string
}

type UserRepositoryI interface {
	CheckExistence(email string) error
	Register(user *User) error
	GetUser(user *User) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(email string) error
	ForgotPassword(data ResetTokenS) error
	GetTokenData(email string) (ResetTokenS, error)
	DeleteTokenData(email string) error
	UpdatePassword(email, password string) error
	GetUserByEmail(email string) (*User, error)
}

type UserUsecaseI interface {
	RegisterUsecase(user *User) error
	VerifyOTPUsecase(user *User) error
	LoginUsecase(user *User) (string, error)
	ForgotPasswordUsecase(email string) error
	ResetPasswordUsecase(data ResetTokenS) error
	OauthCallbackUsecase(user *goth.User) (string, error)
	GetUserByEmail(email string) (*User, error)
}

type PasswordServiceI interface {
	HashPassword(password string) ([]byte, error)
	Compare(hashed, newP string) bool
}

type MailerI interface {
	SendOTPEmail(toEmail, otp string) error
	SendResetPassEmail(toEmail, token string) error
}

type JwtServI interface {
	CreateToken(user User) (string, error)
}

type GeneratorI interface {
	GenerateOTP() string
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}
