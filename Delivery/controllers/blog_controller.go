package controllers

import (
	usecases "blog_api/Usecases"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	UseCase usecases.BlogUseCaseI
}

func NewBlogController(Uc usecases.BlogUseCaseI) *BlogController{
	return &BlogController{
		UseCase: Uc,
	}
}

func (BlgCtrl *BlogController) UpdateBlogController(c *gin.Context) {}