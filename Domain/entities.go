package Domain

import (
	"time"
)

type User struct {
	Username string
	Email    string
	Password string
	Bio      string
	Role     string
	Verfied  bool
	OTP      string
	OTPTime  time.Time
	Provider string
}

type Blog struct {
	ID          string
	Title       string
	Content     string
	Owner_email string
	Tags        []string
	Date        time.Time
	ViewCount   int
	Comments    []string
}

type ResetTokenS struct {
	Email       string
	Token       string
	Created_at  time.Time
	NewPassword string
}

type LikeTracker struct {
	BlogID    string
	UserEmail string
	Liked     int
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply *string `json:"reply"`
}
