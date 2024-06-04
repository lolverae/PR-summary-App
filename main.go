package main

import (
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GithubToken   string
	GmailUsername string
	GmailToken    string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment file")
	}

	config := &Config{
		GithubToken:   os.Getenv("GITHUB_TOKEN"),
		GmailToken:    os.Getenv("GMAIL_TOKEN"),
		GmailUsername: os.Getenv("GMAIL_USERNAME"),
	}

	return config, nil
}

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

func main() {
  report := "Test Report"
  SendReport(report)

}
