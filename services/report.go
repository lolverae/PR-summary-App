package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GenerateReport() (string) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	client, err := createGitHubClient(config.GithubToken)
	if err != nil {
		log.Panicf("Error creating GitHub client: %s", err)
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
		log.Panicf("Error listing pull requests: %s", err)
	}

	open := countState(repos, "open")
	closed := countState(repos, "closed")
	inProgress := open - closed

  summary := fmt.Sprintf(" Hello team,\n Here's the summary of pull requests activity in the last week for the %s repository: \n\nPull Request Summary:\nOpened: %d\nClosed: %d\nIn Progress: %d\n\n", githubRepo, open, closed, inProgress)

	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	var summaryListBuilder strings.Builder

	appendPR := func(pr *github.PullRequest, state string) {
		var stateLabel string
		switch state {
		case "open":
			stateLabel = "Opened"
		case "closed":
			stateLabel = "Closed"
		}

		updatedAt := pr.CreatedAt
		if state == "closed" {
			updatedAt = pr.ClosedAt
		}

		if updatedAt.After(oneWeekAgo) {
			fmt.Fprintf(&summaryListBuilder, "#%d: \"%s\" by %s %s on %s\n", *pr.Number, *pr.Title, *pr.User.Login, stateLabel, updatedAt.Format("January 2, 2006"))
		}
	}

	summaryListBuilder.WriteString("\n- Opened Pull Requests:\n")
	for _, pr := range repos {
		if *pr.State == "open" {
			appendPR(pr, "open")
		}
	}

	summaryListBuilder.WriteString("\n- Closed Pull Requests:\n")
	for _, pr := range repos {
		if *pr.State == "closed" {
			appendPR(pr, "closed")
		}
	}

	emailEnd := "\nPlease review and take necessary actions.\nBest regards,\nYour Name"
	return summary + summaryListBuilder.String() + emailEnd
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
