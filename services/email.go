package services

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendReport(report string) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	mail := composeEmail(report, config)

	err = sendEmail(mail, config)
	if err != nil {
		log.Printf("Failed to send email: %s", err)
		return
	}

	log.Print("Email sent successfully")
}

func composeEmail(report string, config *Config) string {
	return fmt.Sprintf("From: %s\nTo: %s\nSubject: Weekly Pull Request Summary\n\n%s", config.GmailUsername, config.GmailUsername, report)
}

func sendEmail(mail string, config *Config) error {
	auth := smtp.PlainAuth("", config.GmailUsername, config.GmailToken, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, config.GmailUsername, []string{config.GmailUsername}, []byte(mail))
	return err
}
