package utils

import (
	"log"

	"Go_Backend/config"
	"github.com/go-mail/mail"
)

type EmailData struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
}

func SendEmail(data EmailData) error {
	cfg := config.LoadConfig()

	m := mail.NewMessage()
	m.SetHeader("From", cfg.Email)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)

	if data.Text != "" {
		m.SetBody("text/plain", data.Text)
	}

	if data.HTML != "" {
		m.AddAlternative("text/html", data.HTML)
	}

	d := mail.NewDialer("smtp.gmail.com", 465, cfg.Email, cfg.Pass)
	d.SSL = true

	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", data.To, err)
		return err
	}

	log.Printf("✅ Email sent successfully to %s", data.To)
	return nil
}
