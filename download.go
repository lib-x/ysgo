package ysgo

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

const defaultFileDomain = "ysepan.com"

// DownloadURLOptions controls download URL generation behavior.
type DownloadURLOptions struct {
	ForceDownload bool
	TextTimestamp string
	FileDomain    string
}

// BuildDownloadURL constructs a direct download URL for a remote file entry.
func (c *YSClient) BuildDownloadURL(directoryNumber int, downloadToken string, file RemoteFile, opts *DownloadURLOptions) (string, error) {
	if file.Server == "" {
		return "", fmt.Errorf("build download url: empty server")
	}
	if downloadToken == "" {
		return "", fmt.Errorf("build download url: empty download token")
	}
	if file.FileToken == "" {
		return "", fmt.Errorf("build download url: empty file token")
	}
	if file.FileName == "" {
		return "", fmt.Errorf("build download url: empty file name")
	}

	domain := defaultFileDomain
	forceDownload := false
	textTimestamp := ""
	if opts != nil {
		if opts.FileDomain != "" {
			domain = opts.FileDomain
		}
		forceDownload = opts.ForceDownload
		textTimestamp = opts.TextTimestamp
	}

	host := downloadHost(file.Server, domain)
	userEndpoint := c.userEndpoint
	if directoryNumber == 100 {
		userEndpoint = "ys168cs1"
	}

	token := downloadToken
	if forceDownload {
		token = "_" + token
	}

	u := fmt.Sprintf(
		"https://%s/wap/%s/%s/%s/%s",
		host,
		url.PathEscape(userEndpoint),
		url.PathEscape(token),
		url.PathEscape(file.FileToken),
		url.PathEscape(file.FileName),
	)

	if strings.EqualFold(filepath.Ext(file.FileName), ".txt") && textTimestamp != "" {
		q := url.Values{}
		q.Set("", textTimestamp)
		encoded := strings.TrimPrefix(q.Encode(), "=")
		if forceDownload {
			u += "?" + encoded + "&lx=xz"
		} else {
			u += "?" + encoded
		}
	}

	return u, nil
}

func downloadHost(server, domain string) string {
	if strings.EqualFold(server, "X") {
		return "y.ys168.com:8000"
	}
	return fmt.Sprintf("ys-%s.%s", strings.ToLower(server), domain)
}
