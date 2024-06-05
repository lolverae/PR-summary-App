package services

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GithubToken   string
	GmailUsername string
	GmailToken    string
	TargetRepo    string
	RepoOwner     string
	TargetEmails  string
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
		TargetRepo:    os.Getenv("TARGET_REPOSITORY"),
		RepoOwner:     os.Getenv("REPOSITORY_OWNER"),
		TargetEmails:  os.Getenv("TARGET_EMAILS"),
	}

	return config, nil
}
