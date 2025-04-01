package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/CodexVeritax/extractum/pkg/model"
)

type Client struct {
	httpClient         *http.Client
	baseURL            string
	token              string
	userAgent          string
	rateLimitRemaining int
	rateLimitReset     time.Time
}

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func NewClient(token string, options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   "http://api.github.com",
		token:     token,
		userAgent: "extractum/1.0",
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func ParseRepoUrl(repoUrl string) (owner, repo string, err error) {
	re := regexp.MustCompile(`github\.com[/:]([^/]+)/([^/]+?)(?:\.git)?$`)
	matches := re.FindStringSubmatch(repoUrl)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("Invalid Github repository URL : %s", repoUrl)
	}

	return matches[1], matches[2], nil
}

// FetchRepository fetches repository information docs : https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
func (c *Client) FetchRepository(ctx context.Context, owner, repo string) (*model.Repository, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s", owner, repo)

	var repository model.Repository
	if err := c.get(ctx, endpoint, nil, &repository); err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	return &repository, nil

}

func (c *Client) FetchIssues(ctx context.Context, owner, repo string, options map[string]string) ([]model.Issue, error) {

	endpoint := fmt.Sprintf("/repos/%s/%s/issues", owner, repo)

	params := make(url.Values)
	if state, ok := options["state"]; ok {
		params.Set("state", state)
	} else {
		params.Set("state", "all")
	}

	var allIssues []model.Issue
	page := 1
	perPage := 100

	for {
		params.Set("page", fmt.Sprint("%d", page))
		params.Set("per_page", fmt.Sprintf("%d", perPage))

		var issues []model.Issue
		if err := c.get(ctx, endpoint, params, &issues); err != nil {
			return nil, fmt.Errorf("failed to fetch issues: %w", err)
		}

		if len(issues) == 0 {
			break
		}

		// filter out pull requests
		for _, issue := range issues {
			if issue.PullRequest == nil {
				allIssues = append(allIssues, issue)
			}
		}

		if len(issues) < perPage {
			break
		}
		page++

	}
	return allIssues, nil
}

// FetchPullRequests fetches pull requests from a repository
func (c *Client) FetchPullRequests(ctx context.Context, owner, repo string, options map[string]string) ([]model.PullRequest, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)

	params := make(url.Values)
	if state, ok := options["state"]; ok {
		params.Set("state", state)
	} else {
		params.Set("state", "closed")
	}

	var allPRs []model.PullRequest
	page := 1
	perPage := 100

	for {
		params.Set("page", fmt.Sprintf("%d", page))
		params.Set("per_page", fmt.Sprintf("%d", perPage))

		var prs []model.PullRequest
		if err := c.get(ctx, endpoint, params, &prs); err != nil {
			return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
		}

		if len(prs) == 0 {
			break
		}

		// Filter merged PRs if needed
		if options["merged"] == "true" {
			for _, pr := range prs {
				if !pr.MergedAt.IsZero() {
					allPRs = append(allPRs, pr)
				}
			}
		} else {
			allPRs = append(allPRs, prs...)
		}

		if len(prs) < perPage {
			break
		}

		page++
	}

	return allPRs, nil
}

// FetchIssueComments fetches comments for an issue
func (c *Client) FetchIssueComments(ctx context.Context, owner, repo string, issueNumber int) ([]model.Comment, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, issueNumber)

	var allComments []model.Comment
	page := 1
	perPage := 100

	for {
		params := url.Values{}
		params.Set("page", fmt.Sprintf("%d", page))
		params.Set("per_page", fmt.Sprintf("%d", perPage))

		var comments []model.Comment
		if err := c.get(ctx, endpoint, params, &comments); err != nil {
			return nil, fmt.Errorf("failed to fetch comments: %w", err)
		}

		if len(comments) == 0 {
			break
		}

		allComments = append(allComments, comments...)

		if len(comments) < perPage {
			break
		}

		page++
	}

	return allComments, nil
}

func (c *Client) get(ctx context.Context, endpoint string, params url.Values, result interface{}) error {
	apiURL := c.baseURL + endpoint

	if params != nil {
		apiURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	if c.rateLimitRemaining < 5 && time.Now().Before(c.rateLimitReset) {
		waitTime := time.Until(c.rateLimitReset) + time.Second
		fmt.Printf("Rate limit almost exceeded , Waiting for %v\n", waitTime)

		select {
		case <-time.After(waitTime):

		case <-ctx.Done():
			return ctx.Err()
		}
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		c.rateLimitRemaining, _ = strconv.Atoi(remaining)
	}

	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		resetTime, _ := strconv.ParseInt(reset, 10, 64)
		c.rateLimitReset = time.Unix(resetTime, 0)
	}

	if resp.StatusCode == http.StatusForbidden && c.rateLimitRemaining == 0 {
		waitTime := time.Until(c.rateLimitReset) + time.Second

		fmt.Printf("Rate limit exceeded . Waiting for %v\n", waitTime)

		select {
		case <-time.After(waitTime):
			return c.get(ctx, endpoint, params, result)

		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Github API failed: %d - %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to parse the respose: %w", err)
	}

	return nil

}
