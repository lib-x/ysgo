package ysgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *YSClient) newRequest(method, path string, form url.Values) (*http.Request, error) {
	return c.newRequestWithContext(context.Background(), method, path, form)
}

func (c *YSClient) newRequestWithContext(ctx context.Context, method, path string, form url.Values) (*http.Request, error) {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, c.apiBaseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create request %s %s: %w", method, path, err)
	}

	req.Header.Set("Authorization", c.generateAuthToken().String())
	req.Header.Set("DNT", "1")
	req.Header.Set("Pragma", "no-cache")
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func (c *YSClient) do(req *http.Request, out any, successStatus ...int) error {
	resp, err := c.getHTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("execute request %s %s: %w", req.Method, req.URL.String(), err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body %s %s: %w", req.Method, req.URL.String(), err)
	}

	if !isSuccessStatus(resp.StatusCode, successStatus...) {
		return parseYSError(resp.StatusCode, body, http.StatusText(resp.StatusCode))
	}

	if out == nil {
		return nil
	}

	switch target := out.(type) {
	case *[]byte:
		*target = bytes.Clone(body)
		return nil
	default:
		if len(body) == 0 {
			return nil
		}
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("decode response %s %s: %w", req.Method, req.URL.String(), err)
		}
		return nil
	}
}

func isSuccessStatus(code int, allowed ...int) bool {
	if len(allowed) == 0 {
		return code >= http.StatusOK && code < http.StatusMultipleChoices
	}
	for _, status := range allowed {
		if code == status {
			return true
		}
	}
	return false
}

func joinCSV(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, ",")
}

func joinSubdirectories(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, "//<")
}
