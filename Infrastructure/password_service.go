package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func (ps PasswordService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (ps PasswordService) Compare(hashed, newP string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(newP)) == nil
}