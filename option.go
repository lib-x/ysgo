package ysgo

import (
	"net/http"
	"time"
)

// ClientOption configures a YSClient.
type ClientOption func(*YSClient)

// WithAPIBaseURL sets the API base URL.
func WithAPIBaseURL(apiBaseURL string) ClientOption {
	return func(c *YSClient) {
		c.apiBaseURL = apiBaseURL
	}
}

// WithApiBaseUrl is a backward-compatible alias for WithAPIBaseURL.
func WithApiBaseUrl(apiBaseURL string) ClientOption {
	return WithAPIBaseURL(apiBaseURL)
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *YSClient) {
		if client != nil {
			c.c = client
		}
	}
}

// WithAuthToken sets a fixed auth token.
func WithAuthToken(token string) ClientOption {
	return func(c *YSClient) {
		if token != "" {
			c.authToken = token
		}
	}
}

// WithAllowedUploadHosts overrides allowed upload host suffixes.
func WithAllowedUploadHosts(hosts ...string) ClientOption {
	return func(c *YSClient) {
		if len(hosts) > 0 {
			c.allowedUploadHosts = append([]string(nil), hosts...)
		}
	}
}

// WithManagementDirectory overrides the directory number used during admin login.
func WithManagementDirectory(directoryNumber string) ClientOption {
	return func(c *YSClient) {
		if directoryNumber != "" {
			c.managementDirectory = directoryNumber
		}
	}
}

// WithSpacePassword sets the optional space access password used for visitor-gated spaces.
func WithSpacePassword(password string) ClientOption {
	return func(c *YSClient) {
		c.spacePassword = password
	}
}

// WithTimeout overrides the default HTTP timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *YSClient) {
		if timeout > 0 && c.c != nil {
			c.c.Timeout = timeout
		}
	}
}
