package pkg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found")
	}
}

// SendVerificationEmail sends a verification email
func SendVerificationEmail(email string, token string) error {
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpEmail == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP credentials are missing")
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpEmail)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Verify Your Account")
	mailer.SetBody("text/html", fmt.Sprintf("<h3>Click the link to verify your account:</h3><a href='http://localhost:8080/verify?token=%s'>Verify</a>", token))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, smtpPassword)

	// Send email and check for errors
	if err := dialer.DialAndSend(mailer); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to", email)
	return nil
}
