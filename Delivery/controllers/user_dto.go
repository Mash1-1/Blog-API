package controllers

import "time"

type UserDTO struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Bio      string    `json:"bio"`
	Role     string    `json:"role"`
	Verfied  bool      `json:"verifed"`
	OTP      string    `json:"otp"`
	OTPTime  time.Time `json:"otptime"`
	Provider string    `json:"provider"`
}

type UpdateProfileDTO struct {
	Username string `json:"username,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

type ResetTokenSDTO struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	Created_at  time.Time
	NewPassword string `json:"new_password"`
}

type RoleUpdateDTO struct {
	Role string `json:"role" binding:"required"`
}