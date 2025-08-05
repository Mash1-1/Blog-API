package usecases

import (
	"blog_api/Domain"
	"errors"
	"log"
	"unicode"
)

type UserUsecase struct {
	repo Domain.UserRepositoryI
	pass_serv Domain.PasswordServiceI
}

func NewUserUsecase(r Domain.UserRepositoryI, ps Domain.PasswordServiceI) UserUsecase {
	return UserUsecase{
		repo: r,
		pass_serv: ps,
	}
}

func (uc UserUsecase) RegisterUsecase(user *Domain.User) error {
	// Check if user has valid credentials before moving on to insert into db
	if !isValidEmail(user.Email) || !isValidPassword(user.Password){
		return errors.New("invalid email")
	}
	if uc.repo.CheckExistence(user.Email) {
		return errors.New("email already exists in database")
	}
	log.Print("Password : " + user.Password)
	new_p, err := uc.pass_serv.HashPassword(user.Password)
	if err != nil {
		return err 
	}
	log.Print("Hashed password: " + string(new_p))
	user.Password = string(new_p)
	uc.repo.Register(user)
	return nil
}

// email validation function
func isValidEmail(email string) bool {
	n := len(email)
	number_of_at := 0 
	ind := -1
	for i := 0; i<n; i++ {
		if email[i] == '@' {
			ind = i 
			number_of_at += 1
		}
	}
	if number_of_at != 1 || ind == 0{
		return false 
	}
	hasDot := false 
	for i := ind + 1; i < n; i++ {
		if email[i] == '.' {
			hasDot = true
			if i == n-1 {
				return  false 
			}
		}
	}
	return hasDot
}

// Password validation function
func isValidPassword(pass string) bool {
	upper := 0
	lower := 0
	special := 0
	digit := 0

	for _, char := range(pass) {
		if unicode.IsDigit(char) {
			digit += 1
		}else if unicode.IsUpper(char) {
			upper += 1
		} else if unicode.IsLower(char) {
			lower += 1
		}else {special += 1}
	}
	return (len(pass) >= 8 && upper > 0 && special > 0 && lower > 0 && digit > 0)
} 