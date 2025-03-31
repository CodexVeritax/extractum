package api

import (
	"net/http"
	"time"
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
