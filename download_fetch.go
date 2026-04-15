package ysgo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// DownloadToWriter downloads the file at downloadURL into dst.
func (c *YSClient) DownloadToWriter(downloadURL string, dst io.Writer) error {
	return c.DownloadToWriterContext(context.Background(), downloadURL, dst)
}

// DownloadToWriterContext downloads the file at downloadURL into dst with context control.
func (c *YSClient) DownloadToWriterContext(ctx context.Context, downloadURL string, dst io.Writer) error {
	if downloadURL == "" {
		return fmt.Errorf("download to writer: empty url")
	}
	if dst == nil {
		return fmt.Errorf("download to writer: nil writer")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("download to writer: create request: %w", err)
	}

	resp, err := c.getHTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("download to writer: execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("download to writer: status %d: %s", resp.StatusCode, string(body))
	}

	if _, err := io.Copy(dst, resp.Body); err != nil {
		return fmt.Errorf("download to writer: copy body: %w", err)
	}
	return nil
}

// DownloadBytes downloads the file at downloadURL into memory.
func (c *YSClient) DownloadBytes(downloadURL string) ([]byte, error) {
	return c.DownloadBytesContext(context.Background(), downloadURL)
}

// DownloadBytesContext downloads the file at downloadURL into memory with context control.
func (c *YSClient) DownloadBytesContext(ctx context.Context, downloadURL string) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := c.DownloadToWriterContext(ctx, downloadURL, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
