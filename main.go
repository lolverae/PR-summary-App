package main

import (
    "log"
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
      GithubToken:    os.Getenv("GITHUB_TOKEN"),
      GmailToken:     os.Getenv("GMAIL_TOKEN"),
      GmailUsername:  os.Getenv("GMAIL_USERNAME"),
    }

    return config, nil
}


func main() {

}
