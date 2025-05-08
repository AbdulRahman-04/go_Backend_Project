package utils

import (
	"log"

	"github.com/go-mail/mail"
	 "Go_Backend/config"// Apna module name use karo
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, text, html string) error {
	cfg := config.LoadConfig() // Load config
	m := mail.NewMessage()

	m.SetHeader("From", cfg.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	if text != "" {
		m.SetBody("text/plain", text)
	}
	if html != "" {
		m.SetBody("text/html", html)
	}

	// SMTP Server Configuration
	d := mail.NewDialer("smtp.gmail.com", 465, cfg.Email, cfg.Pass)
	d.SSL = true

	// Sending Email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Printf("Email sent successfully to %s ✅", to)
	return nil
}