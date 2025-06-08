package utils

import (
	"log"

	"Go_Backend/config"
	"github.com/go-mail/mail"
)

// EmailData holds the email details.
type EmailData struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
}

// emailQueue is a buffered channel used to queue email tasks.
var emailQueue chan EmailData

// init initializes the email queue and starts the dispatcher goroutine.
func init() {
	// Adjust capacity as needed.
	emailQueue = make(chan EmailData, 100)
	go emailDispatcher()
}

// emailDispatcher continuously reads from the emailQueue and sends emails.
func emailDispatcher() {
	for emailTask := range emailQueue {
		if err := sendEmailInternal(emailTask); err != nil {
			log.Printf("❌ Failed to send email to %s: %v", emailTask.To, err)
		} else {
			log.Printf("✅ Email sent successfully to %s", emailTask.To)
		}
	}
}

// sendEmailInternal sends an email synchronously using the go-mail library.
func sendEmailInternal(data EmailData) error {
	cfg := config.LoadConfig()

	// Create new email message.
	m := mail.NewMessage()
	// Use provided "From" header, or fallback to config.
	if data.From == "" {
		data.From = cfg.Email
	}
	m.SetHeader("From", data.From)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)

	// Set email body parts.
	if data.Text != "" {
		m.SetBody("text/plain", data.Text)
	}
	if data.HTML != "" {
		m.AddAlternative("text/html", data.HTML)
	}

	// Prepare dialer using SMTP settings from config.
	d := mail.NewDialer("smtp.gmail.com", 465, cfg.Email, cfg.Pass)
	d.SSL = true

	// Send the email synchronously.
	err := d.DialAndSend(m)
	return err
}

// QueueEmail enqueues an email task, so it gets sent asynchronously.
func QueueEmail(data EmailData) {
	emailQueue <- data
}