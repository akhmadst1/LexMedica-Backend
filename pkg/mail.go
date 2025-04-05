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
	frontendUrl := os.Getenv("FRONTEND_URL")

	if smtpEmail == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP credentials are missing")
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpEmail)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "LexMedica Account Verification")
	mailer.SetBody("text/html", fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; border: 1px solid #e0e0e0; padding: 20px; border-radius: 10px;">
        <h2 style="color: #263238;">Selamat Datang di LexMedica 👋</h2>
        <p>Email Anda baru saja didaftarkan untuk akun LexMedica melalui halaman berikut %s. Untuk melanjutkan, mohon verifikasi email berikut.</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s/verify-email?token=%s" style="
                background-color: #263238;
                color: white;
                padding: 12px 24px;
                text-decoration: none;
                font-weight: bold;
                border-radius: 6px;
                display: inline-block;
            ">
                Verifikasi Email
            </a>
        </div>
        <p>Jika tombol tidak berfungsi, silakan akses link berikut:</p>
        <p style="word-break: break-all;"><a href="%s/verify-email?token=%s">%s/verify-email?token=%s</a></p>
        <p style="margin-top: 40px; font-size: 0.9em; color: #888;">Jika Anda tidak mendaftarkan akun ini, mohon segera hubungi lexmedica@gmail.com</p>
    </div>`,
		frontendUrl, frontendUrl, token, frontendUrl, token, frontendUrl, token))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, smtpEmail, smtpPassword)

	// Send email and check for errors
	if err := dialer.DialAndSend(mailer); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to", email)
	return nil
}
