package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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

	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	summaryList := "- Opened Pull Requests:\n"
	for _, pr := range repos {
		if *pr.State == "open" && pr.UpdatedAt.After(oneWeekAgo) {
			summaryList += fmt.Sprintf("#%d: \"%s\" by %s opened on %s\n", *pr.Number, *pr.Title, *pr.User.Login, pr.CreatedAt.Format("January 2, 2006"))
		}
	}

	summaryList += "\n- Closed Pull Requests:\n"
	for _, pr := range repos {
		if *pr.State == "closed" && pr.UpdatedAt.After(oneWeekAgo) {
			summaryList += fmt.Sprintf("#%d: \"%s\" by %s closed on %s\n", *pr.Number, *pr.Title, *pr.User.Login, pr.ClosedAt.Format("January 2, 2006"))
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
