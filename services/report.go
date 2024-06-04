package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GenerateReport() string {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	client, err := createGitHubClient(config.GithubToken)
	if err != nil {
		log.Fatalf("Error creating GitHub client: %s", err)
	}

	githubRepo := "kubernetes"
	repoOwner := "kubernetes"

	repos, _, err := client.PullRequests.List(context.TODO(), repoOwner, githubRepo, nil)
	if err != nil {
		log.Println("Error listing pull requests:", err)
		return ""
	}

	return github.Stringify(repos)
}

func createGitHubClient(token string) (*github.Client, error) {
	if token == "" {
		return nil, errors.New("GitHub token is empty")
	}
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}
