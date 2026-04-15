package ysgo

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	defaultAPIBaseURL = "https://c6.ysepan.com"
	defaultTimeout    = 30 * time.Second
)

type YSClient struct {
	apiBaseURL          string
	userEndpoint        string
	managementPass      string
	spacePassword       string
	authToken           string
	managementDirectory string
	allowedUploadHosts  []string

	mu sync.RWMutex
	c  *http.Client
}

func NewClient(userEndpoint string, managementPass string, opts ...ClientOption) *YSClient {
	client := &YSClient{
		userEndpoint:        userEndpoint,
		managementPass:      managementPass,
		authToken:           newAuthToken(),
		managementDirectory: defaultDirectoryNumber,
		allowedUploadHosts:  append([]string(nil), allowedUploadHosts...),
		c:                   &http.Client{Timeout: defaultTimeout},
	}

	applyDefaultOption(client)

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func applyDefaultOption(client *YSClient) {
	if client.apiBaseURL == "" {
		client.apiBaseURL = defaultAPIBaseURL
	}
}

func (c *YSClient) GetAPIBaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiBaseURL
}

// GetApiBaseURL is a backward-compatible alias for GetAPIBaseURL.
func (c *YSClient) GetApiBaseURL() string {
	return c.GetAPIBaseURL()
}

// GetApiBaseUrl is a backward-compatible alias for GetAPIBaseURL.
func (c *YSClient) GetApiBaseUrl() string {
	return c.GetAPIBaseURL()
}

func (c *YSClient) GetUserEndpoint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userEndpoint
}

func (c *YSClient) GetAuthToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authToken
}

func (c *YSClient) getSpacePassword() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.spacePassword
}

func (c *YSClient) getManagementDirectory() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.managementDirectory
}

func (c *YSClient) setAuthToken(token string) {
	if token == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.authToken = token
}

func (c *YSClient) getHTTPClient() *http.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.c
}

func (c *YSClient) getAllowedUploadHosts() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]string(nil), c.allowedUploadHosts...)
}

func newAuthToken() string {
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli(), rand.Intn(10000))
}
