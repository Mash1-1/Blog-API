package Domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
)

type BlogRepositoryI interface {
	Create(blog *Blog) error
	UpdateBlog(updatedBlog *Blog) error
	GetAllBlogs(limit int, offset int) ([]Blog, error)
	SearchBlog(searchBlog *Blog) ([]Blog, error)
	DeleteBlog(id string) error
	FilterBlog(filterBlog *Blog) ([]Blog, error)
	GetBlog(id string) (Blog, error)
	FindLiked(user_email, blog_id string) (*LikeTracker, error)
	CreateLikeTk(lt LikeTracker) error
	DeleteLikeTk(lt LikeTracker) error
	NumberOfDislikes(id string) (int64, error)
	NumberOfLikes(id string) (int64, error)
	GetPopularBlogs() ([]Blog, error)
}

type BlogUseCaseI interface {
	CreateBlogUC(Blog) error
	UpdateBlogUC(Blog) error
	GetAllBlogUC(limit int, offset int) ([]Blog, error)
	SearchBlogUC(Blog) ([]Blog, error)
	DeleteBlogUC(string) error
	FilterBlogUC(Blog) ([]Blog, error)
	GetByIdBlogUC(string) (Blog, error)
	AIChatBlogUC(ChatRequest) (*string, error)
	CheckIfLiked(user_email, blogId string) (int, error)
	AddLikeUC(LikeTracker) error
	Dislikes(id string) (int64, error)
	Likes(id string) (int64, error)
	GetPopularBlogs() ([]Blog, error)
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
	StoreToken(RefreshTokenStorage) error
	GetRefreshToken(string) (string, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUserProfile(user *User) (*User, error)
	UpdateUserRole(email string, role string) (*User, error)
	DeleteToken(email string) error
	AddToBlackList(token, email string) error
	DeleteFromBlackList(email string) error
	RetriveFromBlackList(email string) (string, error)
}

type UserUsecaseI interface {
	RegisterUsecase(user *User) error
	VerifyOTPUsecase(user *User) error
	LoginUsecase(user *User) (map[string]string, error)
	ForgotPasswordUsecase(email string) error
	ResetPasswordUsecase(data ResetTokenS) error
	OauthCallbackUsecase(user *goth.User) (string, error)
	GetUserByEmail(email string) (*User, error)
	UpdateProfileUsecase(user *User) (*User, error)
	UpdateUserRole(email string, role string) (*User, error)
	RefreshUseCase(refreshToken string) (map[string]string, error)
	LogoutUseCase(email, accessToken string) error
	RetriveFromBlackList(email string) (string, error)
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
	CreateToken(user User) (map[string]string, error)
	ParseToken(string) (*jwt.Token, error)
	IsExpired(*jwt.Token) bool
}

type GeneratorI interface {
	GenerateOTP() string
}
type RefreshTokenStorage struct {
	Email string
	Token string
}
