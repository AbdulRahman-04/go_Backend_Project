package utils

import (
	"log"

	"github.com/go-mail/mail"
	"Go_Backend/config"
)

// EmailData struct for email sending
type EmailData struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
}

// SendEmail sends an email using SMTP
func SendEmail(data EmailData) error {
	cfg := config.LoadConfig()
	m := mail.NewMessage()

	m.SetHeader("From", data.From)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)
	if data.Text != "" {
		m.SetBody("text/plain", data.Text)
	}
	if data.HTML != "" {
		m.SetBody("text/html", data.HTML)
	}

	// SMTP Server Configuration
	d := mail.NewDialer("smtp.gmail.com", 465, cfg.Email, cfg.Pass)
	d.SSL = true

	// Sending Email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Printf("Email sent successfully to %s ✅", data.To)
	return nil
}