package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

	repos := listPullRequests(client, repoOwner, githubRepo)
	open, closed := countPullRequestStates(repos)
	inProgress := open - closed

	openPullRequests, closedPullRequests := extractPullRequests(repos)

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

func countPullRequestStates(repos []*github.PullRequest) (int, int) {
	open, closed := 0, 0
	for _, pr := range repos {
		if *pr.State == "open" {
			open++
		} else if *pr.State == "closed" {
			closed++
		}
	}
	return open, closed
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

func listPullRequests(client *github.Client, repoOwner, githubRepo string) []*github.PullRequest {
	opts := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
	}

	repos, _, err := client.PullRequests.List(context.TODO(), repoOwner, githubRepo, opts)
	if err != nil {
		log.Panicf("Error listing pull requests: %s", err)
	}

	return repos
}

func extractPullRequests(repos []*github.PullRequest) ([]PullRequest, []PullRequest) {
	var openPullRequests, closedPullRequests []PullRequest
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	for _, pr := range repos {
		pullRequest := PullRequest{
			Number: *pr.Number,
			Title:  *pr.Title,
			Author: *pr.User.Login,
			Date:   pr.GetCreatedAt().Format("January 2, 2006"),
			URL:    *pr.HTMLURL,
		}

		if pr.GetUpdatedAt().After(oneWeekAgo) {
			if *pr.State == "open" {
				openPullRequests = append(openPullRequests, pullRequest)
			} else if *pr.State == "closed" {
				closedPullRequests = append(closedPullRequests, pullRequest)
			}
		}
	}
	return openPullRequests, closedPullRequests
}
