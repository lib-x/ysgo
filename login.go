package ysgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	loginPath              = "/nkj/csxx.aspx"
	logoutPath             = "/nkj/csxx.aspx"
	loginCommand           = "yzglmm"
	logoutCommand          = "tcgl"
	defaultDirectoryNumber = "1445856"
)

func (c *YSClient) Login() (*LoginResponse, error) {
	reqURL := c.apiBaseUrl + loginPath + "?cz=" + loginCommand

	form := url.Values{}
	form.Add("glmm", c.managementPass)
	form.Add("mlbh", defaultDirectoryNumber)

	req, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token := c.generateAuthToken()
	req.Header.Set("Authorization", token.String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("DNT", "1")
	req.Header.Set("Pragma", "no-cache")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewYSError(resp.StatusCode, "Login failed", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	return &loginResp, nil
}

func (c *YSClient) Logout() error {
	reqURL := c.apiBaseUrl + logoutPath + "?cz=" + logoutCommand

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	token := c.generateAuthToken()
	req.Header.Set("Authorization", token.String())
	req.Header.Set("DNT", "1")
	req.Header.Set("Pragma", "no-cache")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute logout request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewYSError(resp.StatusCode, "Logout failed", resp.Status)
	}

	return nil
}
