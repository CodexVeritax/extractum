package model

import "time"

// Repository represents a GitHub repository , further schema entries found at : https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
type Repository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	OpenIssues  int       `json:"open_issues"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	URL         string    `json:"url"`
}

// Issue represents a GitHub issue , complete schema found at : https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28
type Issue struct {
	ID        int64      `json:"id"`
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	URL       string     `json:"url"`
	Labels    []Label    `json:"labels"`
	Author    User       `json:"author"`

	// Parsed data
	CodeBlocks    []CodeBlock `json:"code_blocks,omitempty"`
	ErrorMessages []string    `json:"error_messages,omitempty"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID        int64      `json:"id"`
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at,omitempty"`
	URL       string     `json:"url"`
	Author    User       `json:"author"`
	Labels    []Label    `json:"labels"`
}

// Comment represents a GitHub comment
type Comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    User      `json:"author"`
	URL       string    `json:"url"`

	// Parsed data
	CodeBlocks    []CodeBlock `json:"code_blocks,omitempty"`
	ErrorMessages []string    `json:"error_messages,omitempty"`
}

// User represents a GitHub user
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

// Label represents a GitHub issue/PR label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CodeBlock represents a code snippet extracted from text
type CodeBlock struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// PatternMatch represents a match against a problem pattern
type PatternMatch struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	URL        string    `json:"url,omitempty"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
}

// IssueCluster represents a group of similar issues
type IssueCluster struct {
	Representative Issue   `json:"representative"`
	Issues         []Issue `json:"issues"`
	Count          int     `json:"count"`
}

// AnalysisResult represents the final analysis output
type AnalysisResult struct {
	Repository    Repository                `json:"repository"`
	Stats         Stats                     `json:"stats"`
	Patterns      map[string][]PatternMatch `json:"patterns"`
	IssueClusters []IssueCluster            `json:"issue_clusters"`
}

// Stats represents repository statistics
type Stats struct {
	TotalIssues int `json:"total_issues"`
	MergedPRs   int `json:"merged_prs"`
}

// IssueDetails contains an issue with all its comments
type IssueDetails struct {
	Issue    Issue     `json:"issue"`
	Comments []Comment `json:"comments"`
}
