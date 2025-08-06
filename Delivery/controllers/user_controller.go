package controllers

import (
	"blog_api/Domain"
	"net/http"
	"time"

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
	Verfied  bool `json:"verifed"`
	OTP  	string `json:"otp"`
	OTPTime time.Time `json:"otptime"`
}

func NewUserController(uc Domain.UserUsecaseI) UserController {
	return UserController{
		usecase: uc,
	}
}

func (UsrCtrl *UserController) LoginController(c *gin.Context) {
	var user UserDTO 
	if c.ShouldBindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid json format"})
		return 
	} 

	token, err := UsrCtrl.usecase.LoginUsecase(UsrCtrl.ChangeToDomain(user))
	if err != nil {
		if err.Error() == "invalid password or email" {
			c.JSON(http.StatusUnauthorized, gin.H{"error" : err.Error()})
			return 
		}
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"message" : "user not found"})
			return 
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return 
	}
	c.JSON(http.StatusOK, gin.H{"message" : "logged in successfully", "token" : token})
}

func (UsrCtrl *UserController) RegisterController(c *gin.Context) {
	var user UserDTO
	err := c.ShouldBind(&user)
	// Handle Binding Errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return 
	}

	err = UsrCtrl.usecase.RegisterUsecase(UsrCtrl.ChangeToDomain(user))
	
	// Handle invalid requests
	if err != nil && (err.Error() == "invalid email" || err.Error() == "email already exists in database"){
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return 
	}

	// Handle DB failure
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return 
	}

	// Handle OTP verification
	c.JSON(http.StatusOK, gin.H{"message" : "OTP sent to your email", "redirect" : "/verify-otp"})
}

func (UsrCtrl *UserController) VerifyOTPController(c *gin.Context) {
	var user UserDTO 

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid user format"})
		return
	}

	user.OTPTime = time.Now()
	err = UsrCtrl.usecase.VerifyOTPUsecase(UsrCtrl.ChangeToDomain(user))
	
	if err != nil && err.Error() == "expired otp code" {
		c.JSON(http.StatusGone, gin.H{"error" : err.Error()})
		return 
	}
	if err != nil && err.Error() == "user not found" {
		c.JSON(http.StatusNotFound, gin.H{"error" : "no user with such email"})
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
		Verfied: user.Verfied,
		OTP: user.OTP,
		OTPTime: user.OTPTime,
	}
	return &dom_user 
}