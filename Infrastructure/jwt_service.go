package infrastructure

import (
	"blog_api/Domain"

	"github.com/dgrijalva/jwt-go"
)

type Jwt_serv struct{}

var JwtSecret = []byte("blog api is amazing")

func (js Jwt_serv) CreateToken(user Domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role" : user.Role,
		"email" : user.Email,
	})
	return token.SignedString(JwtSecret)
}