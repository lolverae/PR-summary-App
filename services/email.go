package services

import (
	"log"
	"net/smtp"
)

func SendReport(report string) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	mail := "From: " + config.GmailUsername + "\n" +
		"To: " + config.GmailUsername + "\n" +
		"Subject: Hello there\n\n" +
		report

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", config.GmailUsername, config.GmailToken, "smtp.gmail.com"),
		config.GmailUsername,
		[]string{config.GmailUsername},
		[]byte(mail),
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}
	log.Print("Sent")
}
