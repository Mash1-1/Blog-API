package usecases

import (
	"blog_api/Domain"
	"errors"
	"log"
	"time"
	"unicode"
)

type UserUsecase struct {
	repo Domain.UserRepositoryI
	pass_serv Domain.PasswordServiceI
	mailer Domain.MailerI
	otpGen Domain.GeneratorI
	jwtServ Domain.JwtServI
}

func NewUserUsecase(r Domain.UserRepositoryI, ps Domain.PasswordServiceI, mailr Domain.MailerI, og Domain.GeneratorI, jt Domain.JwtServI) UserUsecase {
	return UserUsecase{
		repo: r,
		pass_serv: ps,
		mailer: mailr,
		otpGen: og,
		jwtServ: jt,
	}
}

func (uc UserUsecase) ResetPasswordUsecase(data Domain.ResetTokenS) error {
	if !isValidEmail(data.Email) {
		return errors.New("invalid email")
	}
	
	if uc.repo.CheckExistence(data.Email) != nil {
		return errors.New("user not found")
	}

	existingData, err := uc.repo.GetTokenData(data.Email)
	if err != nil {
		return err 
	}
	// Check token expiry and validity
	if data.Created_at.Sub(existingData.Created_at).Minutes() > 10 || data.Token != existingData.Token {
		uc.repo.DeleteTokenData(data.Email)
		return errors.New("invalid token")
	}

	// validate and update password
	if !isValidPassword(data.NewPassword) {
		return errors.New("invalid password")
	}
	hashed, err := uc.pass_serv.HashPassword(data.NewPassword)
	if err != nil {
		return err 
	}
	return uc.repo.UpdatePassword(data.Email, string(hashed))
}

func (uc UserUsecase) ForgotPasswordUsecase(email string) error {
	if !isValidEmail(email) {
		return errors.New("invalid email")
	}

	if uc.repo.CheckExistence(email) != nil {
		return errors.New("user not found")
	}
	
	reset_token := uc.otpGen.GenerateOTP()
	err := uc.mailer.SendOTPEmail(email, reset_token)
	if err != nil {
		return err 
	}
	data := Domain.ResetTokenS{Token: reset_token, Email: email, Created_at: time.Now()}
	return uc.repo.ForgotPassword(data)
}

func (uc UserUsecase) RegisterUsecase(user *Domain.User) error {
	// Check if user has valid credentials before moving on to insert into db
	if !isValidEmail(user.Email) || !isValidPassword(user.Password){
		return errors.New("invalid email")
	}
	if uc.repo.CheckExistence(user.Email) == nil {
		return errors.New("email already exists in database")
	}
	new_p, err := uc.pass_serv.HashPassword(user.Password)
	if err != nil {
		return err 
	}
	if !user.Verfied {
		otp := uc.otpGen.GenerateOTP()
		user.OTP = otp 
		user.OTPTime = time.Now()
		err = uc.mailer.SendOTPEmail(user.Email, otp)
		if err != nil {
			log.Print(err.Error())
			return  errors.New("error while sending otp email")
		}
	}
	user.Password = string(new_p)
	return uc.repo.Register(user)
}

func (uc UserUsecase) LoginUsecase(user *Domain.User) (string, error) {
	existingUser, err := uc.repo.GetUser(user)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !uc.pass_serv.Compare(existingUser.Password, user.Password) {
		return "", errors.New("invalid password or email")
	}
	// Get token using jwt
	return uc.jwtServ.CreateToken(*user)
}

func (uc UserUsecase) VerifyOTPUsecase(user *Domain.User) error {
	existingUser, err := uc.repo.GetUser(user)
	if err != nil {
		return errors.New("user not found")
	}
	if (user.OTPTime.Sub(existingUser.OTPTime)).Minutes() > 5 {
		err := uc.repo.DeleteUser(existingUser.Email)
		if err != nil {
			return err
		}
		return errors.New("expired otp code please register again")
	}
	if (user.OTP != existingUser.OTP) {
		err := uc.repo.DeleteUser(existingUser.Email)
		if err != nil {
			return err
		}
		return  errors.New("invalid otp code, please register again")
	}
	existingUser.Verfied = true
	return uc.repo.UpdateUser(existingUser)
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