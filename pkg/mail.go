package pkg

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// SendVerificationEmail sends a verification email
func SendVerificationEmail(email string, token string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "your-email@example.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Verify Your Account")
	mailer.SetBody("text/html", fmt.Sprintf("<h3>Click the link to verify your account:</h3><a href='http://localhost:8080/verify?token=%s'>Verify</a>", token))

	dialer := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-email-password")

	return dialer.DialAndSend(mailer)
}
