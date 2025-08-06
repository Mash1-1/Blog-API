package infrastructure

import (
	"fmt"
	"math/rand"
	"time"
)

type Generator struct{}

// Generate OTP
func (g Generator) GenerateOTP() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	otp := r.Intn(899999) + 100000

	return fmt.Sprintf("%06d", otp)
}