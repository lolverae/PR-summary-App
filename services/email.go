package services

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"text/template"
	"github.com/k3a/html2text"
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

	log.Print(html2text.HTML2Text(mail.String()))
	log.Print("Email sent successfully")
}

func composeEmail(report EmailData) bytes.Buffer {
	t, err := template.ParseFiles("templates/template.html")
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

	recipients := strings.Split(config.TargetEmails, ",")
	senderAddr := config.GmailUsername

	err := smtp.SendMail("smtp.gmail.com:587", auth, senderAddr, recipients, mail.Bytes())
	return err
}
