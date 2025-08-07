package infrastructure

import (
	"blog_api/Domain"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Jwt_serv struct{}

var JwtSecret = []byte("blog api is amazing")

func (js Jwt_serv) CreateToken(user Domain.User) (map[string]string, error) {
	tokens := make(map[string]string)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Minute * 1),
		"role":  user.Role,
		"email": user.Email,
	})

	t, err := token.SignedString(JwtSecret)
	if err != nil {
		return tokens, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["exp"] = time.Now().Add(time.Hour * 24)
	rtClaims["sub"] = user.Email
	rtClaims["iat"] = time.Now()

	rt, err := refreshToken.SignedString(JwtSecret)
	if err != nil {
		return tokens, err
	}
	tokens["access_token"] = t
	tokens["refresh_token"] = rt
	return tokens, nil
}

func validateToken(tokenString string, user Domain.User, c *gin.Context) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		tkClaims := token.Claims.(jwt.MapClaims)
		if time.Since(tkClaims["exp"].(time.Time)) > 0 {
			return nil, fmt.Errorf("Token Expired")
		}
		c.Set("role", user.Role)
		c.Set("user", user)
		return JwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return token, nil
}
