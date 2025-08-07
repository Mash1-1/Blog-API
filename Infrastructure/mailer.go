package infrastructure

import (
	"net/smtp"
)

type Mailer struct{
	smtpHost string 
	smtpPort string
	smtpUsername string
	smtpPass  string
	from  string
}

func NewMailer(Host, Port, Username, Pass, frm string) Mailer {
	return Mailer{
		smtpHost: Host,
		smtpPort: Port,
		smtpUsername: Username,
		smtpPass: Pass,
		from: frm,
	}
}

func (m *Mailer) SendOTPEmail(toEmail, otp string) error {
	to := []string{toEmail}
	addr := m.smtpHost + ":" + m.smtpPort 

	// Create the message
	subject := "Subject : Your OTP Code: \n"
	body := "Your OTP is: " + otp + " it expires in 5 minutes."
	msg := []byte(subject + body)

	auth := smtp.PlainAuth("", m.smtpUsername, m.smtpPass, m.smtpHost)
	return smtp.SendMail(addr, auth, m.from, to, msg)
}