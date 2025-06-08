package utils

import (
	"log"

	"Go_Backend/config"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSData holds the details for an SMS message.
type SMSData struct {
	To   string
	Body string
}

// smsQueue is a buffered channel to queue SMS tasks.
var smsQueue chan SMSData

// init initializes the SMS queue and starts the dispatcher.
func init() {
	// Adjust the capacity as per expected load.
	smsQueue = make(chan SMSData, 100)
	go smsDispatcher()
}

// smsDispatcher continuously listens to the smsQueue and sends SMS.
func smsDispatcher() {
	for smsTask := range smsQueue {
		if err := sendSMSInternal(smsTask); err != nil {
			log.Printf("❌ Failed to send SMS to %s: %v", smsTask.To, err)
		} else {
			log.Printf("✅ SMS sent successfully to %s", smsTask.To)
		}
	}
}

// sendSMSInternal sends an SMS synchronously using Twilio.
func sendSMSInternal(data SMSData) error {
	cfg := config.LoadConfig()
	log.Printf("SID: %s, Token: %s, Phone: %s", cfg.SID, cfg.Token, cfg.Phone)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.SID,
		Password: cfg.Token,
	})

	params := &api.CreateMessageParams{}
	params.SetTo(data.To)
	params.SetFrom(cfg.Phone)
	params.SetBody(data.Body)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	log.Printf("✅ SMS sent successfully to %s | Message SID: %s", data.To, *resp.Sid)
	return nil
}

// QueueSMS enqueues an SMS task to be sent asynchronously.
func QueueSMS(to, body string) {
	smsQueue <- SMSData{
		To:   to,
		Body: body,
	}
}