package infrastructure

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type Mailer struct{
	smtpHost string 
	smtpPort int
	smtpUsername string
	smtpPass  string
	from  string
}

func NewMailer(Host string, Port int, Username, Pass, frm string) Mailer {
	return Mailer{
		smtpHost: Host,
		smtpPort: Port,
		smtpUsername: Username,
		smtpPass: Pass,
		from: frm,
	}
}

func (m *Mailer) SendOTPEmail(toEmail, otp string) error {
	message := gomail.NewMessage()

	// Compose message
	message.SetHeader("From", m.from)
	message.SetHeader("To", toEmail)
	message.SetHeader("Subject", "Verify your email with this OTP.")
	message.SetBody("text/plain", fmt.Sprintf("This is your OTP : %v", otp))
	
	// Set up the smtp dialer
	dialer := gomail.NewDialer(m.smtpHost, m.smtpPort, m.smtpUsername, m.smtpPass)
	return dialer.DialAndSend(message)
}

func (m *Mailer) SendResetPassEmail(toEmail, token string) error {
	message := gomail.NewMessage()

	// Compose message
	message.SetHeader("From", m.from)
	message.SetHeader("To", toEmail)
	message.SetHeader("Subject", fmt.Sprintf("Reset your email using this token, it expires in 10 Minutes.%v", token))

	// Set up the smtp dialer 
	dialer := gomail.NewDialer(m.smtpHost, m.smtpPort, m.smtpUsername, m.smtpPass)
	return dialer.DialAndSend(message)
}