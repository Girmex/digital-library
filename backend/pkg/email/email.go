package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type Config struct {
	From     string
	Password string
	SmtpHost string
	SmtpPort string
}

type Service struct {
	config Config
}

func NewService(cfg Config) *Service {
	return &Service{config: cfg}
}

func (s *Service) SendVerificationEmail(to, token string) error {
	// Email content
	subject := "Subject: Verify Your Email\r\n"
	body := fmt.Sprintf("Click this link to verify your account: http://localhost:8080/auth/verify-email?token=%s\r\n", token)
	msg := []byte(subject + "\r\n" + body)

	// SMTP auth
	auth := smtp.PlainAuth("", s.config.From, s.config.Password, s.config.SmtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SmtpHost, s.config.SmtpPort)
	fmt.Println(addr)
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Verification email sent to %s", to)
	return nil
}

func (s *Service) SendResetPasswordEmail(to, token string) error {
	// Email content
	subject := "Subject: Reset Your Password\r\n"
	body := fmt.Sprintf("Click this link to reset your password: http://localhost:8080/auth/reset-password?token=%s\r\n", token)
	msg := []byte(subject + "\r\n" + body)

	// SMTP auth
	auth := smtp.PlainAuth("", s.config.From, s.config.Password, s.config.SmtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SmtpHost, s.config.SmtpPort)
	fmt.Println(addr)
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Reset password email sent to %s", to)
	return nil
}
