package services

import (
	"context"
	"errors"
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

	client, err := createGitHubClient(config.GithubToken)
	if err != nil {
		log.Fatalf("Error creating GitHub client: %s", err)
	}

	githubRepo := "kubernetes"
	repoOwner := "kubernetes"

	opts := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
	}

	repos, _, err := client.PullRequests.List(context.TODO(), repoOwner, githubRepo, opts)
	if err != nil {
		log.Println("Error listing pull requests:", err)
		return ""
	}

	open := countState(repos, "open")
	closed := countState(repos, "closed")
	inProgress := open - closed

	summary := fmt.Sprintf("Pull Request Summary:\nOpened: %d\nClosed: %d\nIn Progress: %d\n", open, closed, inProgress)

	summaryList := "- Opened Pull Requests:\n"
	for i, pr := range repos {
		if *pr.State == "open" {
			summaryList += fmt.Sprintf("  %d. #%d: \"%s\" by %s\n", i+1, *pr.Number, *pr.Title, *pr.User.Login)
		}
	}

	summaryList += "\n- Closed Pull Requests:\n"
	for i, pr := range repos {
		if *pr.State == "closed" {
			summaryList += fmt.Sprintf("  %d. #%d: \"%s\" by %s\n", i+1, *pr.Number, *pr.Title, *pr.User.Login)
		}
	}

	summaryList += "\n- In-Progress Pull Requests:\n"
	for i, pr := range repos {
		if *pr.State == "open" && pr.MergedAt == nil {
			summaryList += fmt.Sprintf("  %d. #%d: \"%s\" by %s\n", i+1, *pr.Number, *pr.Title, *pr.User.Login)
		}
	}

	return summary + summaryList
}

func countState(repos []*github.PullRequest, state string) int {
	count := 0
	for _, pr := range repos {
		if *pr.State == state {
			count++
		}
	}
	return count
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
