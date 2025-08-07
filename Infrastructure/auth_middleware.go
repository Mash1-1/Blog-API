package infrastructure

import (
	"blog_api/Domain"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Usecase Domain.UserUsecaseI
}

func (am AuthMiddleware) Auth_role() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if ok && role == "admin" {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you need to be an admin to view this site"})
	}
}

func (am AuthMiddleware) Auth_token() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization header is required", "message": "you need to be logged in to view this site"})
			c.Abort()
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(authParts[0], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return JwtSecret, nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)
		// Check if the JWT is valid and has the type MapClaims
		if err == nil && ok && token.Valid {
			// Get role and store it for the next handlers to authorize role
			c.Set("role", claims["role"].(string))
		} else {
			c.JSON(401, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		user, err := am.Usecase.GetUserByEmail(claims["email"].(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		c.Set("user", user)

		c.Next()
	}
}
