package utils

import (
 "log"

 "github.com/twilio/twilio-go"
 api "github.com/twilio/twilio-go/rest/api/v2010"
 "Go_Backend/config" // Tumhara module name: ensure that "Go_Backend" is correct
)

// SendSMS sends a message using the Twilio API.
func SendSMS(to, body string) error {
 // Load configuration (SID, Token, and Twilio phone number)
 cfg := config.LoadConfig()
 log.Printf("SID: %s, Token: %s, Phone: %s", cfg.SID, cfg.Token, cfg.Phone)

 // Create a Twilio client using credentials.
 client := twilio.NewRestClientWithParams(twilio.ClientParams{
 	Username: cfg.SID,
 	Password: cfg.Token,
 })

 // Create parameters for the SMS.
 params := &api.CreateMessageParams{}
 params.SetTo(to)
 params.SetFrom(cfg.Phone)
 params.SetBody(body)

 // Send the message.
 resp, err := client.Api.CreateMessage(params)
 if err != nil {
 	log.Printf("Failed to send SMS: %v", err)
 	return err
 }

 log.Printf("SMS sent successfully to %s ✅ Message SID: %s", to, *resp.Sid)
 return nil
}