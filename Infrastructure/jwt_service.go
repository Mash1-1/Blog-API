package infrastructure

import (
	"blog_api/Domain"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Jwt_serv struct{}

var JwtSecret = []byte("blog api is amazing")

func (js Jwt_serv) CreateToken(user Domain.User) (map[string]string, error) {
	tokens := make(map[string]string)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"role":  user.Role,
		"email": user.Email,
		"type":  "access",
	})

	t, err := token.SignedString(JwtSecret)
	if err != nil {
		return tokens, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	rtClaims["sub"] = user.Email
	rtClaims["iat"] = time.Now()
	rtClaims["email"] = user.Email
	rtClaims["type"] = "refresh"

	rt, err := refreshToken.SignedString(JwtSecret)
	if err != nil {
		return tokens, err
	}
	tokens["access_token"] = t
	tokens["refresh_token"] = rt
	return tokens, nil
}

func (js Jwt_serv) ParseToken(token string) (*jwt.Token, error) {
	ptoken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return JwtSecret, nil
	})
	return ptoken, err
}

func (js Jwt_serv) IsExpired(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)

	exp, ok := claims["exp"].(float64)
	if !ok {
		return true // No exp claim → treat as expired
	}
	// Convert exp to time.Time and compare
	return time.Unix(int64(exp), 0).Before(time.Now())
}
