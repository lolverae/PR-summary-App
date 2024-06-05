package services

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

func SendReport(report EmailData) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	mail := composeEmail(report)

	err = sendEmail(mail, config)
	if err != nil {
		log.Printf("Failed to send email: %s", err)
		return
	}

	log.Print("Email sent successfully")
}

func composeEmail(report EmailData) bytes.Buffer {
	t, err := template.ParseFiles("services/template.html")
	if err != nil {
		log.Fatalf("Error parsing template: %s", err)
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Weekly Pull Request Summary\n%s\n\n", mimeHeaders)))

	if err := t.Execute(&body, report); err != nil {
		log.Fatalf("Error executing template: %s", err)
	}
	return body
}

func sendEmail(mail bytes.Buffer, config *Config) error {
	auth := smtp.PlainAuth("", config.GmailUsername, config.GmailToken, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, config.GmailUsername, []string{config.GmailUsername}, mail.Bytes())
	return err
}
