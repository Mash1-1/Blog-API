package usecases

import (
	"blog_api/Domain"
	"errors"
	"log"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
)

type UserUsecase struct {
	repo      Domain.UserRepositoryI
	pass_serv Domain.PasswordServiceI
	mailer    Domain.MailerI
	otpGen    Domain.GeneratorI
	jwtServ   Domain.JwtServI
}

func NewUserUsecase(r Domain.UserRepositoryI, ps Domain.PasswordServiceI, mailr Domain.MailerI, og Domain.GeneratorI, jt Domain.JwtServI) UserUsecase {
	return UserUsecase{
		repo:      r,
		pass_serv: ps,
		mailer:    mailr,
		otpGen:    og,
		jwtServ:   jt,
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

func (uc UserUsecase) OauthCallbackUsecase(user *goth.User) (string, error) {

	if uc.repo.CheckExistence(user.Email) == nil {
		// Handle login since user already exists
		existingUser, err := uc.repo.GetUserByEmail(user.Email)
		if err != nil {
			return "", err
		}
		// Get token using jwt
		tokens, err := uc.jwtServ.CreateToken(*existingUser)
		if err != nil {
			return "", err
		}
		return tokens["access_token"], nil
	} else {
		// Register the user into the database and login the user
		var newUser Domain.User
		newUser.Email = user.Email
		err := uc.repo.Register(&newUser)
		if err != nil {
			return "", err
		}
		tokens, err := uc.jwtServ.CreateToken(newUser)
		if err != nil {
			return "", err
		}
		return tokens["access_token"], nil
	}
}

func (uc UserUsecase) ForgotPasswordUsecase(email string) error {
	if !isValidEmail(email) {
		return errors.New("invalid email")
	}

	if uc.repo.CheckExistence(email) != nil {
		return errors.New("user not found")
	}

	reset_token := uc.otpGen.GenerateOTP()
	err := uc.mailer.SendResetPassEmail(email, reset_token)
	if err != nil {
		return err
	}
	data := Domain.ResetTokenS{Token: reset_token, Email: email, Created_at: time.Now()}
	return uc.repo.ForgotPassword(data)
}

func (uc UserUsecase) RegisterUsecase(user *Domain.User) error {
	// Check if user has valid credentials before moving on to insert into db
	if !isValidEmail(user.Email) || !isValidPassword(user.Password) {
		return errors.New("invalid email or password")
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
			return errors.New("error while sending otp email")
		}
	}
	user.Password = string(new_p)
	return uc.repo.Register(user)
}

func (uc UserUsecase) LoginUsecase(user *Domain.User) (map[string]string, error) {
	var tokens map[string]string
	existingUser, err := uc.repo.GetUser(user)
	if err != nil {
		return tokens, errors.New("user not found")
	}

	if !uc.pass_serv.Compare(existingUser.Password, user.Password) {
		return tokens, errors.New("invalid password or email")
	}
	tokens, err = uc.jwtServ.CreateToken(*user)
	if err != nil {
		return tokens, err
	}
	tokenData := Domain.RefreshTokenStorage{
		Email: user.Email,
		Token: tokens["refresh_token"],
	}
	if err = uc.repo.StoreToken(tokenData); err != nil {
		return tokens, err
	}
	// Get token using jwt
	return tokens, nil
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
	if user.OTP != existingUser.OTP {
		err := uc.repo.DeleteUser(existingUser.Email)
		if err != nil {
			return err
		}
		return errors.New("invalid otp code, please register again")
	}
	existingUser.Verfied = true
	return uc.repo.UpdateUser(existingUser)
}

func (uc UserUsecase) GetUserByEmail(email string) (*Domain.User, error) {
	user, err := uc.repo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// email validation function
func isValidEmail(email string) bool {
	n := len(email)
	number_of_at := 0
	ind := -1
	for i := 0; i < n; i++ {
		if email[i] == '@' {
			ind = i
			number_of_at += 1
		}
	}
	if number_of_at != 1 || ind == 0 {
		return false
	}
	hasDot := false
	for i := ind + 1; i < n; i++ {
		if email[i] == '.' {
			hasDot = true
			if i == n-1 {
				return false
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

	for _, char := range pass {
		if unicode.IsDigit(char) {
			digit += 1
		} else if unicode.IsUpper(char) {
			upper += 1
		} else if unicode.IsLower(char) {
			lower += 1
		} else {
			special += 1
		}
	}
	return (len(pass) >= 8 && upper > 0 && special > 0 && lower > 0 && digit > 0)
}

func (uc UserUsecase) RefreshUseCase(refreshToken string) (map[string]string, error) {
	tokens := make(map[string]string)
	token, err := uc.jwtServ.ParseToken(refreshToken)
	if err != nil {
		return tokens, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userEmail := claims["email"].(string)
	if uc.jwtServ.IsExpired(token) {
		return tokens, errors.New("refresh Token expired try to login again")
	}
	if err := uc.repo.DeleteToken(userEmail); err != nil {
		return tokens, err
	}

	user, err := uc.repo.GetUserByEmail(userEmail)
	if err != nil {
		return tokens, err
	}
	tokens, err = uc.jwtServ.CreateToken(*user)
	if err != nil {
		return tokens, err
	}
	if err = uc.repo.StoreToken(Domain.RefreshTokenStorage{Token: tokens["refresh_token"], Email: userEmail}); err != nil {
		return tokens, err
	}
	return tokens, err
}
