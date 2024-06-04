package services

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

	repos, _, err := client.PullRequests.List(ctx, repoOwner, githubRepo, nil)
	if err != nil {
		fmt.Println(err)
	}

	return github.Stringify(repos)
}
