package main

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

func GenerateReport() string {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	githubRepo := "kubernetes"
	repoOwner := "kubernetes"

	repos, _, err := client.PullRequests.List(ctx, repoOwner, githubRepo,
		nil)
	if err != nil {
		fmt.Println(err)
	}

	return github.Stringify(repos)
}

func main() {
	report := GenerateReport()
	SendReport(report)

}
