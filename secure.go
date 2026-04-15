package ysgo

import (
	"fmt"
	"net/url"
	"strings"
)

var allowedUploadHosts = []string{
	"ysepan.com",
	"ys168.com",
}

func validateUploadAddress(raw string, allowed []string) error {
	if raw == "" {
		return fmt.Errorf("empty upload address")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse upload address: %w", err)
	}
	host := strings.ToLower(u.Hostname())
	if u.Scheme != "https" && !isLoopbackHost(host) {
		return fmt.Errorf("upload address must use https: %s", raw)
	}
	for _, suffix := range allowed {
		suffix = strings.ToLower(strings.TrimSpace(suffix))
		if suffix == "" {
			continue
		}
		if host == suffix || strings.HasSuffix(host, "."+suffix) {
			return nil
		}
	}
	if isLoopbackHost(host) {
		return nil
	}
	return fmt.Errorf("upload address host not allowed: %s", host)
}

func isLoopbackHost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
