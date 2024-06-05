package services

type PullRequest struct {
	Number int
	Title  string
	Author string
	Date   string
	URL    string
}

type EmailData struct {
	Repo               string
	Opened             int
	Closed             int
	InProgress         int
	OpenPullRequests   []PullRequest
	ClosedPullRequests []PullRequest
}

