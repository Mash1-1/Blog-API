package controllers

import (
	"blog_api/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase Domain.UserUsecaseI
}

type UserDTO struct {
	Username string `json:"username"`
	Email    string	`json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Role     string `json:"role"`
}

func NewUserController(uc Domain.UserUsecaseI) UserController {
	return UserController{
		usecase: uc,
	}
}

func (UsrCtrl *UserController) RegisterController(c *gin.Context) {
	var user UserDTO
	err := c.ShouldBind(&user)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return 
	}

	err = UsrCtrl.usecase.RegisterUsecase(UsrCtrl.ChangeToDomain(user))
	if err != nil && (err.Error() == "invalid email" || err.Error() == "email already exists in database"){
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return 
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return 
	}
	c.JSON(http.StatusCreated, gin.H{"message" : "user registered successfully."})
}

func (UsrCtrl *UserController) ChangeToDomain(user UserDTO) *Domain.User {
	var dom_user = Domain.User{
		Email: user.Email,
		Password: user.Password,
		Username: user.Username,
		Bio: user.Bio,
		Role: user.Role,
	}
	return &dom_user 
}