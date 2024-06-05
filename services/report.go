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

type PullRequest struct {
	Number int
	Title  string
	Author string
	Date   string
}

type EmailData struct {
	Repo               string
	Opened             int
	Closed             int
	InProgress         int
	OpenPullRequests   []PullRequest
	ClosedPullRequests []PullRequest
}

var openPullRequests []PullRequest
var closedPullRequests []PullRequest

func GenerateReport() EmailData {
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
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

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

	for _, pr := range repos {
		if *pr.State == "open" {
			appendPR(pr, "open")
			openPullRequests = append(openPullRequests, PullRequest{
				Number: *pr.Number,
				Title:  *pr.Title,
				Author: *pr.User.Login,
				Date:   pr.CreatedAt.Format("January 2, 2006"),
			})
		} else if *pr.State == "closed" {
			appendPR(pr, "closed")
			closedPullRequests = append(closedPullRequests, PullRequest{
				Number: *pr.Number,
				Title:  *pr.Title,
				Author: *pr.User.Login,
				Date:   pr.ClosedAt.Format("January 2, 2006"),
			})
		}
	}

	emailData := EmailData{
		Repo:               githubRepo,
		Opened:             open,
		Closed:             closed,
		InProgress:         inProgress,
		OpenPullRequests:   openPullRequests,
		ClosedPullRequests: closedPullRequests,
	}

	return emailData
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
