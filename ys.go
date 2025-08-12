package ysgo

import "net/http"

const (
	defaultApiBaseUrl = "http://c6.ysepan.com"
)

type YSClient struct {
	apiBaseUrl     string
	userEndpoint   string
	managementPass string

	c *http.Client
}

func NewClient(userEndpoint string, managementPass string, opts ...ClientOption) *YSClient {
	client := &YSClient{
		userEndpoint:   userEndpoint,
		managementPass: managementPass,
		c:              &http.Client{},
	}

	applyDefaultOption(client)

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func applyDefaultOption(client *YSClient) {
	if client.apiBaseUrl == "" {
		client.apiBaseUrl = defaultApiBaseUrl
	}
}

func (c *YSClient) GetApiBaseURL() string {
	return c.apiBaseUrl
}

func (c *YSClient) GetUserEndpoint() string {
	return c.userEndpoint
}
